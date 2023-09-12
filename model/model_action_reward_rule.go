package model

import (
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type ActionRewardRule struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	ActionName      string    `json:"action_name"`
	DisplayName     string    `json:"display_name"`
	Amount          float64   `json:"amount"`
	Limit           uint      `json:"limit"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedAt       time.Time `json:"created_at"`
	EntityBalanceId uint      `json:"entity_balance_id"`
}

func (ActionRewardRule) TableName() string {
	return "action_reward_rules"
}

func (m *ActionRewardRule) GetById(id uint) error {
	err := DB().Where("id=?", id).First(m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (m *ActionRewardRule) GetByActionNameAndEntityId(action_name string, entity_balance_id uint) error {
	err := DB().
		Where("action_name=?", action_name).
		Where("entity_balance_id=?", entity_balance_id).
		First(m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (m *ActionRewardRule) Create() error {
	db := DB().Create(m)

	var mysqlErr *mysql.MySQLError
	if errors.As(db.Error, &mysqlErr) && mysqlErr.Number == DbDuplicateEntryCode {
		return ErrKeyConflict
	} else if db.Error != nil {
		return db.Error
	}

	return nil
}

func (m *ActionRewardRule) Update() error {
	err := DB().Model(&m).Save(m).Error

	return err
}

func (m *ActionRewardRule) Delete() error {
	db := DB().Delete(ActionRewardRule{}, "id = ?", m.ID)

	if db.Error != nil {
		return db.Error
	}

	return nil
}

func GetAllSubworldRewardRules(entityBalanceId uint) ([]ActionRewardRule, error) {
	var rs []ActionRewardRule
	err := DB().
		Where("entity_balance_id = ?", entityBalanceId).
		Find(&rs).Error
	return rs, err
}
