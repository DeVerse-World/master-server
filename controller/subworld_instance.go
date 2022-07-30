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

type SubworldInstanceController struct {
}

func NewSubworldInstanceController() *SubworldInstanceController {
	return &SubworldInstanceController{}
}

func (ctrl *SubworldInstanceController) GetAll(c *gin.Context) {
	const (
		success = "Get All Subworld Instances successfully"
		failed  = "Get All Subworld Instances unsuccessfully"
	)

	userIdStr := c.Request.URL.Query().Get("user_id")

	if userIdStr == "" {
		sis, err := model.GetAllSubworldInstances()
		if err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}
		JSONReturn(c, http.StatusOK, success, gin.H{
			"subworld_instances": sis,
		})
	} else {
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}
		sis, err := model.GetSubworldInstancesFromHost(userId)
		if err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}
		JSONReturn(c, http.StatusOK, success, gin.H{
			"subworld_instances": sis,
		})
	}
}

func (ctrl *SubworldInstanceController) Create(c *gin.Context) {
	const (
		success = "Create Subworld Instance successfully"
		failed  = "Create Subworld Instance unsuccessfully"
	)

	var req requestSchema.CreateSubworldInstance
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

	var subworld_instance = req.SubworldInstance
	subworld_instance.HostId = &wallet.ID
	subworld_instance.NumCurrentPlayers = 0
	if err := subworld_instance.Create(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"subworld_instance": subworld_instance,
	})
}

func (ctrl *SubworldInstanceController) Update(c *gin.Context) {
	const (
		success = "Update Subworld Instance successfully"
		failed  = "Update Subworld Instance unsuccessfully"
	)

	var req requestSchema.UpdateSubworldInstance
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

	var subworld_instance model.SubworldInstance
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if err := subworld_instance.GetById(id); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if *subworld_instance.HostId != wallet.ID {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("unauthorized to update other's subworld instance"))
		return
	}

	subworld_instance.HostName = req.HostName
	subworld_instance.Region = req.Region
	subworld_instance.MaxPlayers = req.MaxPlayers
	subworld_instance.NumCurrentPlayers = req.NumCurrentPlayers
	subworld_instance.InstancePort = req.InstancePort
	subworld_instance.BeaconPort = req.BeaconPort
	if err := subworld_instance.Update(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"subworld_instance": subworld_instance,
	})
}

func (ctrl *SubworldInstanceController) Delete(c *gin.Context) {
	const (
		success = "Delete Subworld Instance successfully"
		failed  = "Delete Subworld Instance unsuccessfully"
	)

	authW, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authW == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var wallet model.Wallet
	wallet.GetWalletByAddress(authW.Address)

	var subworld_instance model.SubworldInstance
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if err := subworld_instance.GetById(id); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if *subworld_instance.HostId != wallet.ID {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("unauthorized to delete other's subworld instance"))
	}
	if err := subworld_instance.Delete(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{})
}
