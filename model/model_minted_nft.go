package model

import "time"

type MintedNft struct {
	ID                          uint      `gorm:"primary_key" json:"id"`
	TokenAddress                string    `json:"token_address"`
	TokenId                     string    `json:"token_id"`
	Name                        string    `json:"name"`
	Description                 string    `json:"description"`
	Supply                      int       `json:"supply"`
	AssetType                   string    `json:"asset_type"`
	FileAssetName               string    `json:"file_asset_name"`
	FileAssetUri                string    `json:"file_asset_uri"`
	FileAssetUriFromCentralized string    `json:"file_asset_uri_from_centralized"`
	File2dUri                   string    `json:"file_2d_uri"`
	File3dUri                   string    `json:"file_3d_uri"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}

func (MintedNft) TableName() string {
	return "minted_nfts"
}

func (n *MintedNft) Create() error {
	db := DB().Create(n)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}
