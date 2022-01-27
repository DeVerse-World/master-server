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
	w, err := jwt.HandleUserCookie(c.Writer, c.Request)
	if err == nil {
		var wallet model.Wallet
		err2 := wallet.GetWalletByAddress(w.Address)
		if err2 == nil {
			c.JSON(http.StatusOK, wallet)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (ctrl *WalletController) UpdateAssets(c *gin.Context) {
	var req schema.UpdateAssetReq
	c.BindJSON(&req)

	wallet, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if err := wallet.GetWalletByAddress(wallet.Address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// TODO: Create In Batch
	for i := 0; i < len(req.Nfts); i++ {
		var nft model.Nft
		var collection model.NftCollection
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
		var wallet model.Wallet
		wallet.GetWalletByAddress(req.Address)
		// store login Link approval if session key exists
		if req.SessionKey != "" {
			var lr model.LoginRequest
			if err := lr.GetByKey(req.SessionKey); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			} else {
				lr.UpdateWalletId(wallet.ID)
			}
		}
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

func (ctrl *WalletController) CreateLoginLink(c *gin.Context) {

	var lr model.LoginRequest
	if err := lr.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"login_url": lr.SessionKey})
	}
}

func (ctrl *WalletController) AuthLoginLink(c *gin.Context) {
	var req schema.AuthLoginLink
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

			if req.SessionKey != "" {
				var lr model.LoginRequest
				if err := lr.GetByKey(req.SessionKey); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				} else {
					lr.UpdateWalletId(wallet.ID)
				}
			}

			c.JSON(http.StatusOK, wallet)
			wallet.UpdateNonce()
		}
	}
}

func (ctrl *WalletController) PollLoginLink(c *gin.Context) {
	sessionKey := c.Param("session_key")

	var lr model.LoginRequest
	if err := lr.GetByKey(sessionKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if lr.WalletId == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no wallet id => this log in link is not authorized yet"})
		return
	}

	var wallet model.Wallet
	if err := wallet.GetWalletById(*lr.WalletId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt.WriteUserCookie(c.Writer, &wallet)
	c.JSON(http.StatusOK, wallet)
}
