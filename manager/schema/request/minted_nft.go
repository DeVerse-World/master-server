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
