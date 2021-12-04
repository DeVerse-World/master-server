package model

import (
	"errors"
	"github.com/hyperjiang/gin-skeleton/manager"
	"gorm.io/gorm"
	"time"
)

type Wallet struct {
	ID uint `gorm:"primary_key" json:"id"`
	Address string `json:"address"`
	Nonce string `json:"nonce"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Wallet) TableName() string {
	return "wallets"
}

func (w *Wallet) Create() error {
	w.Nonce = manager.Utils{}.GenerateRandomString(10)
	db := DB().Create(w)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (w *Wallet) GetWalletByAddress(address string) error {
	err := DB().Where("address=?", address).First(w).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}
