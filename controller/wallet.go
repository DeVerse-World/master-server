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

func (ctrl *WalletController) GetOrCreateWallet(c *gin.Context) {
	var req schema.SignupWalletReq
	c.BindJSON(&req)

	authW, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authW.Address == req.Address { // nonce message was signed
		return
	}

	var wallet model.Wallet

	if err := wallet.GetOrCreate(req.Address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, wallet)
	}
}

func (ctrl *WalletController) Auth(c *gin.Context) {
	var req schema.AuthWalletReq
	c.BindJSON(&req)

	var wallet model.Wallet
	if err := wallet.GetWalletByAddress(req.Address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		msg := "I am signing my one-time nonce: " + wallet.Nonce
		verifyResult := util.VerifySig(wallet.Address, req.Signature, []byte(msg))

		if !verifyResult {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			jwt.WriteUserCookie(c.Writer, &wallet)
			c.JSON(http.StatusOK, wallet)
			wallet.UpdateNonce()
		}
	}

}

func (ctrl *WalletController) MockAuth(c *gin.Context) {
	var req schema.MockAuthWalletReq
	c.BindJSON(&req)

	var wallet model.Wallet
	wallet.Address = req.Address

	jwt.WriteUserCookie(c.Writer, &wallet)
}
