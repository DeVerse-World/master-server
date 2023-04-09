package model

import (
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type EntityBalance struct {
	ID            uint      `gorm:"primary_key" json:"id"`
	EntityId      int       `json:"entity_id"`
	EntityType    string    `json:"entity_type"`
	BalanceAmount uint      `json:"balance_amount"`
	BalanceType   string    `json:"balance_type"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (EntityBalance) TableName() string {
	return "entity_balances"
}

func (eb *EntityBalance) GetById(id int) error {
	err := DB().Where("id=?", id).First(eb).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (eb *EntityBalance) GetBySubworldIdAndType(subworld_id int, balance_type string) error {
	err := DB().
		Where("entity_type=?", SUBWORLD_TYPE).
		Where("entity_id=?", subworld_id).
		Where("balance_type=?", balance_type).
		First(eb).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (eb *EntityBalance) Create() error {
	db := DB().Create(eb)

	var mysqlErr *mysql.MySQLError
	if errors.As(db.Error, &mysqlErr) && mysqlErr.Number == DbDuplicateEntryCode {
		return ErrKeyConflict
	} else if db.Error != nil {
		return db.Error
	}

	return nil
}

func (eb *EntityBalance) Update() error {
	err := DB().Model(&eb).Save(eb).Error

	return err
}

func (eb *EntityBalance) Delete() error {
	db := DB().Delete(SubworldTemplate{}, "id = ?", eb.ID)

	if db.Error != nil {
		return db.Error
	}

	return nil
}
