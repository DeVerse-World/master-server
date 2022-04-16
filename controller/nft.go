package controller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/hyperjiang/gin-skeleton/manager/schema"
)

type NftController struct{}

func (ctrl *NftController) CreateMintNftLink(c *gin.Context) {
	var req schema.CreateMinkLink
	c.BindJSON(&req)

	c.JSON(http.StatusOK, gin.H{"mint_nft_url": os.Getenv("UI_HOST") + "/marketplace/create-item?fileUri=" + req.IpfsHash})
}
