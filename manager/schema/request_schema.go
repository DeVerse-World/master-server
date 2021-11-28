package schema

type SignupWalletReq struct {
	Address string `json:"address" binding:"required"`
}

type MockAuthWalletReq struct {
	Address string `json:"address" binding:"required"`
}

type AuthWalletReq struct {
	Address string `json:"address" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}
