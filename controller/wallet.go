package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hyperjiang/gin-skeleton/manager/jwt"
	"github.com/hyperjiang/gin-skeleton/manager/schema"
	"github.com/hyperjiang/gin-skeleton/manager/util"
	"github.com/hyperjiang/gin-skeleton/model"
	"net/http"
)

type WalletController struct{}

func (ctrl *WalletController) GetWallet(c *gin.Context) {
	var wallet model.Wallet

	address := c.Param("address")

	if err := wallet.GetWalletByAddress(address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, wallet)
	}
}

func (ctrl *WalletController) GetWalletPrivateProfile(c *gin.Context) {
	w, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	c.JSON(http.StatusOK, w)
}

func (ctrl *WalletController) UpdateAssets(c *gin.Context) {
	var req schema.UpdateAssetReq
	c.BindJSON(&req)

	address := c.Param("address")
	var wallet model.Wallet
	if err := wallet.GetWalletByAddress(address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var nft model.Nft
	var collection model.NftCollection

	// TODO: Create In Batch
	for i := 0; i < len(req.Nfts); i++ {
		collection.TokenAddress = req.Nfts[i].TokenAddress
		collection.GetOrCreate()

		nft.ImageUrl = req.Nfts[i].ImageUrl
		nft.WalletId = wallet.ID
		nft.CollectionId = collection.ID
		nft.Create()
	}
}

func (ctrl *WalletController) FetchAssets(c *gin.Context) {
	address := c.Param("address")
	var wallet model.Wallet
	if err := wallet.GetWalletByAddress(address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if assets, err := wallet.FetchAssetsByAddress(wallet.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, assets)
	}
}

func (ctrl *WalletController) Signup(c *gin.Context) {
	var req schema.SignupWalletReq
	c.BindJSON(&req)

	var wallet model.Wallet
	wallet.Address = req.Address
	wallet.Nonce = util.GenerateRandomString(10)

	if err := wallet.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, wallet)
	}
}

func (ctrl *WalletController) Auth(c *gin.Context) {

}

func (ctrl *WalletController) MockAuth(c *gin.Context) {
	var req schema.MockAuthWalletReq
	c.BindJSON(&req)

	var wallet model.Wallet
	wallet.Address = req.Address

	jwt.WriteUserCookie(c.Writer, &wallet)
}
