package model

import (
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type ActionRewardRecord struct {
	ID                 uint      `gorm:"primary_key" json:"id"`
	OccurCount         uint      `json:"occur_count"`
	UpdatedAt          time.Time `json:"updated_at"`
	CreatedAt          time.Time `json:"created_at"`
	ActionRewardRuleId uint      `json:"action_reward_rule_id"`
	UserId             uint      `json:"user_id"`
}

func (ActionRewardRecord) TableName() string {
	return "action_reward_records"
}

func (m *ActionRewardRecord) GetById(id int) error {
	err := DB().Where("id=?", id).First(m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (m *ActionRewardRecord) GetByRuleAndUser(ruleId uint, userId uint) error {
	err := DB().Where("action_reward_rule_id = ?", ruleId).Where("user_id = ?", userId).First(m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}
	return err
}

func (m *ActionRewardRecord) Create() error {
	db := DB().Create(m)

	var mysqlErr *mysql.MySQLError
	if errors.As(db.Error, &mysqlErr) && mysqlErr.Number == DbDuplicateEntryCode {
		return ErrKeyConflict
	} else if db.Error != nil {
		return db.Error
	}

	return nil
}

func (m *ActionRewardRecord) Update() error {
	err := DB().Model(&m).Save(m).Error

	return err
}

func (m *ActionRewardRecord) Delete() error {
	db := DB().Delete(SubworldTemplate{}, "id = ?", m.ID)

	if db.Error != nil {
		return db.Error
	}

	return nil
}

func GetAllUserSubworldRewardRecords(user_id string, rule_ids []uint) ([]ActionRewardRecord, error) {
	var rs []ActionRewardRecord
	err := DB().
		Where("user_id = ?", user_id).
		Where("action_reward_rule_id IN (?)", rule_ids).
		Find(&rs).Error
	return rs, err
}

func GetAllSubworldRuleRecordsExceedLimit(rule_id uint, newLimit uint) ([]ActionRewardRecord, error) {
	var rs []ActionRewardRecord
	err := DB().
		Where("action_reward_rule_id = ?", rule_id).
		Where("occur_count > ?", newLimit).
		Find(&rs).Error
	if errors.Is(err, gorm.ErrRecordNotFound) || len(rs) == 0 {
		return nil, ErrDataNotFound
	}

	return rs, err
}

func SumAllUserRewardedAmount(user_id uint, balance_type string) (float64, error) {
	var totalValue float64
	err := DB().Table("action_reward_records").
		Joins("JOIN action_reward_rules ON action_reward_records.action_reward_rule_id = action_reward_rules.id").
		Joins("JOIN entity_balances ON action_reward_rules.entity_balance_id = entity_balances.id").
		Where("action_reward_records.user_id = ? AND entity_balances.balance_type = ?", user_id, balance_type).
		Select("COALESCE(SUM(action_reward_records.occur_count * action_reward_rules.amount), 0) as total_value").
		Scan(&totalValue).Error
	return totalValue, err
}
