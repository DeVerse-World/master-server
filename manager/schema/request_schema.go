package schema

type SignupWalletReq struct {
	Address string `json:"address" binding:"required"`
	// UserId string `json:"user_id"`
}

type GetOrCreateWalletReq struct {
	Address string `json:"address" binding:"required"`
}

type MockAuthWalletReq struct {
	Address string `json:"address" binding:"required"`
}

type AuthWalletReq struct {
	Address string `json:"address" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type UpdateAssetReq struct {
	Nfts []NftStruct `json:"nfts"`
}
