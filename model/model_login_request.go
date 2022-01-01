package model

import (
	"errors"
	"github.com/hyperjiang/gin-skeleton/manager/util"
	"gorm.io/gorm"
	"time"
)

type LoginRequest struct {
	ID uint `gorm:"primary_key" json:"id"`
	SessionKey string `json:"session_key"`
	WalletId *uint `json:"wallet_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (LoginRequest) TableName() string {
	return "login_requests"
}

func (lr *LoginRequest) Create() error {
	lr.SessionKey = util.GenerateRandomString(10)
	db := DB().Create(lr)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (lr *LoginRequest) GetByKey(session_key string) error {
	err := DB().Where("session_key=?", session_key).First(lr).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (lr *LoginRequest) UpdateWalletId(wallet_id uint) {
	DB().Model(&lr).Update("wallet_id", wallet_id)
}