package controller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/hyperjiang/gin-skeleton/manager"
	requestSchema "github.com/hyperjiang/gin-skeleton/manager/schema/request"
	"github.com/hyperjiang/gin-skeleton/model"
)

type NftController struct {
	inMemoryStorageManager *manager.InMemoryStorageManager
}

func NewNftController(
	inMemoryStorageManager *manager.InMemoryStorageManager,
) *NftController {
	return &NftController{
		inMemoryStorageManager,
	}
}

func (ctrl *NftController) CreateMintNftLink(c *gin.Context) {
	const (
		success = "Create mint nft link successfully"
		failed  = "Create mint nft link unsuccessfully"
	)

	var req requestSchema.CreateMinkLink
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"mint_nft_url": os.Getenv("UI_HOST") + "/marketplace/create-item?fileUri=" + req.IpfsHash,
	})
}

func (ctrl *NftController) NotifyMinted(c *gin.Context) {
	const (
		success = "Notify Minted successfully"
		failed  = "Notify Minted unsuccessfully"
	)

	var req requestSchema.NotifyMinted
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	var mintedNft = req.MintedNft
	if err := mintedNft.Create(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, mintedNft)
}

func (ctrl *NftController) CheckName(c *gin.Context) {
	const (
		success = "Check Name successfully"
		failed  = "Check Name unsuccessfully"
	)
	name := c.Request.URL.Query().Get("name")

	var data interface{}
	if err := ctrl.inMemoryStorageManager.Get("nft_"+name, &data); err != nil && err.Error() != "cache miss" {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if data != nil {
		JSONReturn(c, http.StatusOK, success, gin.H{"exist": true})
		return
	}

	var mintedNft model.MintedNft
	if err := mintedNft.GetByName(name); err != nil && err != model.ErrDataNotFound {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if mintedNft.ID > 0 {
		JSONReturn(c, http.StatusOK, success, gin.H{"exist": true})
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{"exist": false})
}

func (ctrl *NftController) LockName(c *gin.Context) {
	const (
		success = "Lock Name successfully"
		failed  = "Lock Name unsuccessfully"
	)

	var req requestSchema.LockName
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if err := ctrl.inMemoryStorageManager.Set("nft_"+req.Name, 1); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, nil)
}

func (ctrl *NftController) UnlockName(c *gin.Context) {
	const (
		success = "Unlock Name successfully"
		failed  = "Unlock Name unsuccessfully"
	)

	var req requestSchema.LockName
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if err := ctrl.inMemoryStorageManager.Delete("nft_" + req.Name); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, nil)
}
