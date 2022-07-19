package model

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/hyperjiang/gin-skeleton/manager/util"
)

type Wallet struct {
	ID      uint   `gorm:"primary_key" json:"id"`
	Address string `json:"address"`
	Nonce   string `json:"nonce"`
	// UserId string `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Wallet) TableName() string {
	return "wallets"
}

func (w *Wallet) Create() error {
	w.Nonce = util.GenerateRandomString(10)
	db := DB().Create(w)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (w *Wallet) GetOrCreate(address string) error {
	dbErr := w.GetWalletByAddress(address)

	if dbErr == nil {
		return nil
	} else {
		w.Address = address
		return w.Create()
	}
}

func (w *Wallet) DeleteWalletNft() error {
	db := DB().Delete(Nft{}, "wallet_id = ?", w.ID)

	if db.Error != nil {
		return db.Error
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

func (w *Wallet) GetWalletById(id uint) error {
	err := DB().Where("id=?", id).First(w).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (w *Wallet) UpdateNonce() {
	DB().Model(&w).Update("nonce", util.GenerateRandomString(10))
}

func (w *Wallet) FetchAssetsByAddress(walletID uint) ([]Nft, error) {
	var nfts []Nft
	err := DB().Find(&nfts, "wallet_id = ?", walletID).Error
	return nfts, err
}

func (w *Wallet) GetWalletAvatars(walletID uint) ([]Avatar, error) {
	var avatars []Avatar
	err := DB().Find(&avatars, "wallet_id = ?", walletID).Error
	return avatars, err
}

func GetTemporaryRewards(walletID uint) ([]MintedNft, error) {
	var rewardNfts []MintedNft
	// TODO: only gives if in ranking
	DB().Raw("SELECT mn.* from minted_nfts mn join event_rewards er on er.minted_nft_id = mn.id "+
		"join events e on e.id = er.event_id join event_participants ep on ep.event_id = e.id WHERE e.allow_temporary_hold > 0 "+
		"and ep.wallet_id = ?", walletID).Scan(&rewardNfts)
	return rewardNfts, nil
}
