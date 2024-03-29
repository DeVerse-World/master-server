package model

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/hyperjiang/gin-skeleton/manager/util"
)

type User struct {
	ID            uint      `gorm:"primary_key" json:"id"`
	SocialEmail   string    `json:"social_email"`
	CustomEmail   string    `json:"custom_email"`
	WalletAddress string    `json:"wallet_address"`
	WalletNonce   string    `json:"wallet_nonce"`
	SteamId       string    `json:"steam_id"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

func (w *User) GetOrCreateByWallet(address string) error {
	dbErr := w.GetUserByWalletAddress(address)

	if dbErr == nil {
		return nil
	} else {
		w.WalletAddress = address
		w.WalletNonce = util.GenerateRandomString(10)
		return w.Create()
	}
}

func (w *User) GetOrCreateBySocialEmail(socialEmail string) error {
	dbErr := w.GetUserBySocialEmail(socialEmail)

	if dbErr == nil {
		return nil
	}

	w.SocialEmail = socialEmail
	return w.Create()
}

func (w *User) GetOrCreateBySteamId(steamId string) error {
	dbErr := w.GetUserBySteamId(steamId)

	if dbErr == nil {
		return nil
	} else {
		w.SteamId = steamId
		return w.Create()
	}
}

func (w *User) Create() error {
	db := DB().Create(w)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (w *User) Update() error {
	// This is a tricky hack for now until figure out how to do partial indexing properly in MySQL
	var testUsers []User
	// TODO: only gives if in ranking
	DB().Raw("SELECT * from users WHERE (social_email = ? and social_email != '') or (wallet_address = ? and wallet_address != '')", w.SocialEmail, w.WalletAddress).Scan(&testUsers)
	if len(testUsers) > 1 {
		return errors.New("this wallet address/ social email was linked to another user before")
	}

	err := DB().Model(&w).Save(w).Error

	return err
}

func (w *User) DeleteUserNft() error {
	db := DB().Delete(Nft{}, "user_id = ?", w.ID)

	if db.Error != nil {
		return db.Error
	}

	return nil
}

func (w *User) GetUserByWalletAddress(walletAddress string) error {
	err := DB().Where("wallet_address=?", walletAddress).First(w).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (w *User) GetUserBySocialEmail(socialEmail string) error {
	err := DB().Where("social_email=?", socialEmail).First(w).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (w *User) GetUserByGoogleEmail(googleEmail string) error {
	err := DB().Where("social_email=?", googleEmail).First(w).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (w *User) GetUserBySteamId(steamId string) error {
	err := DB().Where("steam_id=?", steamId).First(w).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (w *User) GetUserById(id string) error {
	err := DB().Where("id=?", id).First(w).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (w *User) GetUserByIdUInt(id uint) error {
	err := DB().Where("id=?", id).First(w).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	}

	return err
}

func (w *User) UpdateWalletNonce() {
	DB().Model(&w).Update("wallet_nonce", util.GenerateRandomString(10))
}

func (w *User) FetchAssetsByUser(userID uint) ([]Nft, error) {
	var nfts []Nft
	err := DB().Find(&nfts, "user_id = ?", userID).Error
	return nfts, err
}

func GetTemporaryRewards(userID uint) ([]MintedNft, error) {
	var rewardNfts []MintedNft
	// TODO: only gives if in ranking
	DB().Raw("SELECT mn.* from minted_nfts mn join event_rewards er on er.minted_nft_id = mn.id "+
		"join events e on e.id = er.event_id join event_participants ep on ep.event_id = e.id WHERE e.allow_temporary_hold > 0 "+
		"and ep.user_id = ?", userID).Scan(&rewardNfts)
	return rewardNfts, nil
}
