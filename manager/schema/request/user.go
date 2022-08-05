package request

import "github.com/hyperjiang/gin-skeleton/manager/schema"

type SignupUserReq struct {
	LoginMode      string `json:"login_mode" binding:"required"`
	WalletAddress  string `json:"wallet_address"`
	GoogleEmail    string `json:"google_email"`
	CustomEmail    string `json:"custom_email"`
	CustomPassword string `json:"custom_password"`
	// UserId string `json:"user_id"`
	SessionKey string `json:"session_key"`
}

type GetOrCreateUserReq struct {
	LoginMode      string `json:"login_mode" binding:"required"`
	WalletAddress  string `json:"wallet_address"`
	GoogleEmail    string `json:"google_email"`
	CustomEmail    string `json:"custom_email"`
	CustomPassword string `json:"custom_password"`
}

type MockAuthUserReq struct {
	LoginMode      string `json:"login_mode" binding:"required"`
	WalletAddress  string `json:"wallet_address"`
	GoogleEmail    string `json:"google_email"`
	CustomEmail    string `json:"custom_email"`
	CustomPassword string `json:"custom_password"`
}

type AuthUserReq struct {
	LoginMode       string `json:"login_mode" binding:"required"`
	WalletAddress   string `json:"address"`
	WalletSignature string `json:"signature"`
	GoogleEmail     string `json:"google_email"`
	GoogleToken     string `json:"google_token"`
	CustomEmail     string `json:"custom_email"`
	CustomPassword  string `json:"custom_password"`
}

type AuthLoginLink struct {
	LoginMode       string `json:"login_mode" binding:"required"`
	SessionKey      string `json:"session_key"`
	WalletAddress   string `json:"wallet_address"`
	WalletSignature string `json:"wallet_signature"`
	GoogleEmail     string `json:"google_email"`
	GoogleToken     string `json:"google_token"`
	CustomEmail     string `json:"custom_email"`
	CustomPassword  string `json:"custom_password"`
}

type UpdateAssetReq struct {
	Nfts []schema.NftStruct `json:"nfts"`
}

type UpdateUserProfileReq struct {
	Name string `json:"name"`
}
