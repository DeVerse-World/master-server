package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	requestSchema "github.com/hyperjiang/gin-skeleton/manager/schema/request"
	"github.com/hyperjiang/gin-skeleton/model"
)

type SystemSettingController struct {
}

func NewSystemSettingController() *SystemSettingController {
	return &SystemSettingController{}
}

func (ctrl *SystemSettingController) GetByInfo(c *gin.Context) {
	const (
		success = "Get System Setting By Info successfully"
		failed  = "Get System Setting By Info unsuccessfully"
	)
	var req requestSchema.GetSystemSettingByInfo
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	var setting model.SystemSetting
	if err := setting.GetByInfo(req.Key, req.Category, req.ObjectReferenceId); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	JSONReturn(c, http.StatusOK, success, gin.H{
		"setting": setting,
	})
}

func (ctrl *SystemSettingController) CreateOrUpdate(c *gin.Context) {
	const (
		success = "Creat Or Update System Setting By Info successfully"
		failed  = "Create Or Update System Setting By Info unsuccessfully"
	)
	var req requestSchema.CreateOrUpdateSystemSetting
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	var systemSetting model.SystemSetting
	systemSetting.KeyName = req.Key
	systemSetting.Category = req.Category
	systemSetting.ObjectReferenceId = req.ObjectReferenceId
	if err := systemSetting.GetByInfo(req.Key, req.Category, req.ObjectReferenceId); err != nil {
		systemSetting.Create()
	}

	systemSetting.Value = req.Value
	if err := systemSetting.Update(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"setting": systemSetting,
	})
}
