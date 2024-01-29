package request

import "github.com/hyperjiang/gin-skeleton/model"

type CreateMinkLink struct {
	IpfsHash string `json:"ipfs_hash" binding:"required"`
}

type NotifyMinted struct {
	MintedNft model.MintedNft `json:"minted_nft" binding:"required"`
}

type LockName struct {
	Name string `json:"name" binding:"required"`
}

type UnlockName struct {
	Name string `json:"name" binding:"required"`
}

type MintNft struct {
	ReceiverIdWalletAddress string `json:"receiver_id_wallet_address"`
	Metadata                struct {
		Static struct {
			Issuer           string `json:"issuer"`
			ContentsProvider string `json:"contents_provider"`
			Name             string `json:"name"`
			Description      string `json:"description"`
			Image            string `json:"image"`
			Link             string `json:"link"`
			SerialNumber     int    `json:"serial_number"`
			TotalIssued      int    `json:"total_issued"`
			AnimationURL     string `json:"animation_url"`
		} `json:"static"`
		Dynamic struct {
		} `json:"dynamic"`
	} `json:"metadata"`
}

type TransferNft struct {
	TokenId                 string `json:"token_id" binding:"required"`
	ReceiverIdWalletAddress string `json:"receiver_id_wallet_address" binding:"required"`
}

type BurnNft struct {
	TokenId string `json:"token_id" binding:"required"`
}

type UpdateDynamicNft struct {
	TokenId         string   `json:"token_id" binding:"required"`
	DynamicMetadata struct{} `json:"dynamic_metadata" binding:"required"`
}
