package model

import "time"

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
