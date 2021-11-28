package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hyperjiang/gin-skeleton/manager/schema"
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

func (ctrl *WalletController) Signup(c *gin.Context) {
	var req schema.SignupWalletReq
	c.BindJSON(&req)

	var wallet model.Wallet
	wallet.Address = req.Address

	if err := wallet.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, wallet)
	}
}

func (ctrl *WalletController) Auth(c *gin.Context) {

}

func (ctrl *WalletController) MockAuth(c *gin.Context) {

}
