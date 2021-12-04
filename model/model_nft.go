package model

import "time"

type Nft struct {
	ID uint `gorm:"primary_key" json:"id"`
	ImageUrl string `json:"image_url"`
	Name string `json:"name"`
	Description string `json:"description"`
	RequireFetch bool `json:"require_fetch"`
	WalletId int `json:"wallet_id"`
	CollectionId int `json:"collection_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Nft) TableName() string {
	return "nfts"
}
