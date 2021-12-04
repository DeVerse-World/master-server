package model

import (
	"time"
)

type NftCollection struct {
	ID uint `gorm:"primary_key" json:"id"`
	TokenAddress string `json:"token_address"`
	Amount string `json:"amount"`
	BlockNumber int `json:"block_number"`
	MintedBlockNum int `json:"minted_block_num"`
	ContractType string `json:"contract_type"`
	TokenUrl string `json:"token_url"`
	TokenId string `json:"token_id"`
	Name string `json:"name"`
	Symbol string `json:"string"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (NftCollection) TableName() string {
	return "nft_collections"
}

func (c *NftCollection) Create() error {
	db := DB().Create(c)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (c *NftCollection) GetOrCreate() (*NftCollection, error) {
	err := DB().Where("token_address=?", c.TokenAddress).First(c).Error
	if err == nil {
		return c, nil
	}
	db := DB().Create(c)

	if db.Error != nil {
		return c, db.Error
	} else if db.RowsAffected == 0 {
		return c, ErrKeyConflict
	}

	return c, nil
}
