package schema

type SignupWalletReq struct {
	Address string `json:"address" binding:"required"`
	// UserId string `json:"user_id"`
	SessionKey string `json:"session_key"`
}

type GetOrCreateWalletReq struct {
	Address string `json:"address" binding:"required"`
}

type MockAuthWalletReq struct {
	Address string `json:"address" binding:"required"`
}

type AuthWalletReq struct {
	Address   string `json:"address" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type AuthLoginLink struct {
	SessionKey string `json:"session_key"`
	Address    string `json:"address" binding:"required"`
	Signature  string `json:"signature" binding:"required"`
}

type CreateMinkLink struct {
	IpfsHash string `json:"ipfs_hash"`
}

type UpdateAssetReq struct {
	Nfts []NftStruct `json:"nfts"`
}
