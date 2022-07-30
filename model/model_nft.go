package model

import "time"

type Nft struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	ImageUrl     string    `json:"image_url"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	RequireFetch bool      `json:"require_fetch"`
	UserId       uint      `json:"user_id"`
	CollectionId uint      `json:"collection_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Nft) TableName() string {
	return "nfts"
}

func (n *Nft) Create() error {
	db := DB().Create(n)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}
