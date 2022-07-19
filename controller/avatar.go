package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/hyperjiang/gin-skeleton/manager/jwt"
	requestSchema "github.com/hyperjiang/gin-skeleton/manager/schema/request"
	"github.com/hyperjiang/gin-skeleton/model"
)

type AvatarController struct{}

func NewAvatarController() *AvatarController {
	return &AvatarController{}
}

func (ctrl *AvatarController) Get(c *gin.Context) {
	const (
		success = "Get Avatar successfully"
		failed  = "Get Avatar unsuccessfully"
	)
	var avatar model.Avatar
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if err := avatar.GetById(id); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"avatar": avatar,
	})
}

func (ctrl *AvatarController) Create(c *gin.Context) {
	const (
		success = "Create Avatar successfully"
		failed  = "Create Avatar unsuccessfully"
	)

	var req requestSchema.CreateAvatar
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	authW, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authW == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var wallet model.Wallet
	wallet.GetWalletByAddress(authW.Address)

	var avatar = req.Avatar
	avatar.WalletId = &wallet.ID
	if err := avatar.Create(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"avatar": avatar,
	})
}

func (ctrl *AvatarController) Update(c *gin.Context) {
	const (
		success = "Update Avatar successfully"
		failed  = "Update Avatar unsuccessfully"
	)

	var req requestSchema.UpdateAvatar
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	authW, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authW == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var wallet model.Wallet
	wallet.GetWalletByAddress(authW.Address)

	var avatar model.Avatar
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
	}
	if err := avatar.GetById(id); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
	}

	if *avatar.WalletId != wallet.ID {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("unauthorized to update other's avatar"))
	}

	if req.PostprocessUrl != "" {
		avatar.PostprocessUrl = req.PostprocessUrl
	}
	avatar.Update()

	JSONReturn(c, http.StatusOK, success, gin.H{
		"avatar": avatar,
	})
}

func (ctrl *AvatarController) Delete(c *gin.Context) {
	const (
		success = "Delete Avatar successfully"
		failed  = "Delete Avatar unsuccessfully"
	)

	authW, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authW == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var wallet model.Wallet
	wallet.GetWalletByAddress(authW.Address)

	var avatar model.Avatar
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
	}
	if err := avatar.GetById(id); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
	}

	if *avatar.WalletId != wallet.ID {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("unauthorized to delete other's avatar"))
	}
	avatar.Delete()

	JSONReturn(c, http.StatusOK, success, gin.H{})
}
