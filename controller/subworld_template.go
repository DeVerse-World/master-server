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

type SubworldTemplateController struct {
}

func NewSubworldTemplateController() *SubworldTemplateController {
	return &SubworldTemplateController{}
}

func (ctrl *SubworldTemplateController) GetAllRoot(c *gin.Context) {
	const (
		success = "Get Root Subworld Templates successfully"
		failed  = "Get Root Subworld Templates unsuccessfully"
	)

	userIdStr := c.Request.URL.Query().Get("user_id")

	if userIdStr == "" {
		sts, err := model.GetAllRoot()
		if err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}
		JSONReturn(c, http.StatusOK, success, gin.H{
			"subworld_templates": sts,
		})
	} else {
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}
		sts, err := model.GetRootFromCreator(userId)
		if err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}
		JSONReturn(c, http.StatusOK, success, gin.H{
			"subworld_templates": sts,
		})
	}
}

func (ctrl *SubworldTemplateController) CreateRoot(c *gin.Context) {
	const (
		success = "Create Root Subworld Template successfully"
		failed  = "Create Root Subworld Template unsuccessfully"
	)

	var req requestSchema.CreateSubworldTemplateRoot
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

	var subworld_template = req.SubworldTemplate
	subworld_template.CreatorId = &wallet.ID
	if err := subworld_template.Create(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"subworld_template": subworld_template,
	})
}

func (ctrl *SubworldTemplateController) UpdateRoot(c *gin.Context) {
	const (
		success = "Update Root Subworld Template successfully"
		failed  = "Update Root Subworld Template unsuccessfully"
	)

	var req requestSchema.UpdateSubworldTemplateRoot
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

	var subworld_template model.SubworldTemplate
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if err := subworld_template.GetById(id); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if *subworld_template.CreatorId != wallet.ID {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("unauthorized to update other's subworld template"))
		return
	}

	subworld_template.FileName = req.FileName
	subworld_template.DisplayName = req.DisplayName
	subworld_template.LevelIpfsUri = req.LevelIpfsUri
	subworld_template.LevelCentralizedUri = req.LevelCentralizedUri
	subworld_template.ThumbnailCentralizedUri = req.ThumbnailCentralizedUri
	if err := subworld_template.Update(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"subworld_template": subworld_template,
	})
}

func (ctrl *SubworldTemplateController) DeleteRoot(c *gin.Context) {
	const (
		success = "Delete Root Subworld Template successfully"
		failed  = "Delete Root Subworld Template unsuccessfully"
	)

	authW, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authW == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var wallet model.Wallet
	wallet.GetWalletByAddress(authW.Address)

	var subworld_template model.SubworldTemplate
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if err := subworld_template.GetById(id); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if *subworld_template.CreatorId != wallet.ID {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("unauthorized to delete other's avatar"))
	}
	if err := subworld_template.Delete(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{})
}
