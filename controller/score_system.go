package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	requestSchema "github.com/hyperjiang/gin-skeleton/manager/schema/request"
	"github.com/hyperjiang/gin-skeleton/model"
)

type ScoreSystemController struct{}

func NewScoreSystemController() *ScoreSystemController {
	return &ScoreSystemController{}
}

func (ctrl *ScoreSystemController) GrantScoreMapping(c *gin.Context) {
	const (
		success = "Grant score mapping successfully"
		failed  = "Grant score mapping unsuccessfully"
	)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	var req requestSchema.GrantScoreMapping
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	model.DB().Transaction(func(tx *gorm.DB) error {
		if err := createEntityBalanceAndActionRules(
			id,
			model.DP_TYPE,
			req.DpMapping,
		); err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return err
		}

		if err := createEntityBalanceAndActionRules(
			id,
			model.EXP_TYPE,
			req.ExpMapping,
		); err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return err
		}

		JSONReturn(c, http.StatusOK, success, gin.H{})
		return nil
	})
}

func createEntityBalanceAndActionRules(
	id int,
	balanceType string,
	balanceMapping requestSchema.BalanceMapping,
) error {
	var dp_entity_balance model.EntityBalance
	dp_entity_balance.EntityId = id
	dp_entity_balance.EntityType = model.SUBWORLD_TYPE
	dp_entity_balance.BalanceAmount = balanceMapping.BalanceAmount
	dp_entity_balance.BalanceType = balanceType

	if err := dp_entity_balance.Create(); err != nil {
		return err
	}

	for key, value := range balanceMapping.ActionRewards {
		var action_reward_rule model.ActionRewardRule
		action_reward_rule.ActionName = key
		action_reward_rule.Amount = value.Amount
		action_reward_rule.Limit = value.Limit
		action_reward_rule.EntityBalanceId = dp_entity_balance.ID
		if err := action_reward_rule.Create(); err != nil {
			return err
		}
	}
	return nil
}

func (ctrl *ScoreSystemController) RetrieveScoreMapping(c *gin.Context) {
	const (
		success = "Retrieve score mapping successfully"
		failed  = "Retrieve score mapping unsuccessfully"
	)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	balance_type := c.Request.URL.Query().Get("type")
	var entity_balance model.EntityBalance
	if err := entity_balance.GetBySubworldIdAndType(id, balance_type); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	action_reward_rules, err := model.GetAllSubworldRewardRules(entity_balance.ID)
	JSONReturn(c, http.StatusOK, success, gin.H{
		"action_reward_rules": action_reward_rules,
	})
}

func (ctrl *ScoreSystemController) UpdateUserScore(c *gin.Context) {
	const (
		success = "Update user score successfully"
		failed  = "Update user score unsuccessfully"
	)

	var req requestSchema.UpdateUserScore
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	model.DB().Transaction(func(tx *gorm.DB) error {
		for _, user_score := range req.Scores {
			for _, user_action := range user_score.DpMapping {
				if err := updateUserScore(
					user_score.UserId,
					user_action.RuleId,
					user_action.RewardedAmount,
				); err != nil {
					abortWithStatusError(c, http.StatusBadRequest, failed, err)
					return err
				}
			}
			for _, user_action := range user_score.ExpMapping {
				if err := updateUserScore(
					user_score.UserId,
					user_action.RuleId,
					user_action.RewardedAmount,
				); err != nil {
					abortWithStatusError(c, http.StatusBadRequest, failed, err)
					return err
				}
			}
		}
		JSONReturn(c, http.StatusOK, success, gin.H{})
		return nil
	})
}

func updateUserScore(
	user_id uint,
	rule_id uint,
	amount uint,
) error {
	var action_reward_rule model.ActionRewardRule
	if err := action_reward_rule.GetById(rule_id); err != nil {
		return err
	}

	if action_reward_rule.Amount*action_reward_rule.Limit < amount {
		return errors.New(
			fmt.Sprintf("user amount exceed limit for user %d and rule %d ",
				user_id, rule_id,
			),
		)
	}

	var action_reward_record model.ActionRewardRecord
	action_reward_record.UserId = user_id
	action_reward_record.ActionRewardRuleId = rule_id
	if err := action_reward_record.GetByRuleAndUser(rule_id, user_id); err == model.ErrDataNotFound {
		if err := action_reward_record.Create(); err != nil {
			return err
		}
	}

	action_reward_record.Amount = amount
	if err := action_reward_record.Update(); err != nil {
		return err
	}

	return nil
}

func (ctrl *ScoreSystemController) RetrieveUserScore(c *gin.Context) {
	const (
		success = "Retrieve score mapping successfully"
		failed  = "Retrieve score mapping unsuccessfully"
	)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	user_id := c.Request.URL.Query().Get("user_id")

	balance_type := c.Request.URL.Query().Get("type")
	var entity_balance model.EntityBalance
	if err := entity_balance.GetBySubworldIdAndType(id, balance_type); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	action_reward_rules, err := model.GetAllSubworldRewardRules(entity_balance.ID)
	if err := entity_balance.GetBySubworldIdAndType(id, balance_type); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	var rule_ids []uint
	for _, rule := range action_reward_rules {
		rule_ids = append(rule_ids, rule.ID)
	}

	action_reward_records, err := model.GetAllUserSubworldRewardRecords(user_id, rule_ids)
	if err := entity_balance.GetBySubworldIdAndType(id, balance_type); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	JSONReturn(c, http.StatusOK, success, gin.H{
		"action_reward_records": action_reward_records,
	})
}
