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

func (ctrl *SubworldTemplateController) GetById(c *gin.Context) {
	const (
		success = "Get Template By Id successfully"
		failed  = "Get Template By Id unsuccessfully"
	)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	var st model.SubworldTemplate
	err = st.GetById(id)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	JSONReturn(c, http.StatusOK, success, gin.H{
		"subworld_template": st,
	})
}

func (ctrl *SubworldTemplateController) GetAllRoot(c *gin.Context) {
	const (
		success = "Get Root Subworld Templates successfully"
		failed  = "Get Root Subworld Templates unsuccessfully"
	)

	var sts []model.SubworldTemplate
	var err error

	userIdStr := c.Request.URL.Query().Get("user_id")
	derivable := c.Request.URL.Query().Get("derivable")

	if userIdStr == "" {
		sts, err = model.GetAllRoot()
		if err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}
	} else {
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}
		sts, err = model.GetRootFromCreator(userId)
		if err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}
	}

	filtered_sts := []model.SubworldTemplate{}
	for i := range sts {
		if derivable == "" || strconv.Itoa(sts[i].Derivable) == derivable {
			filtered_sts = append(filtered_sts, sts[i])
		}
	}
	JSONReturn(c, http.StatusOK, success, gin.H{
		"subworld_templates": filtered_sts,
	})
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

	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authU == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var user model.User
	user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))

	var subworld_template = req.SubworldTemplate
	subworld_template.CreatorId = &user.ID
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

	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authU == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var user model.User
	user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))

	var subworld_template model.SubworldTemplate
	idStr := c.Param("root_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if err := subworld_template.GetById(id); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if *subworld_template.CreatorId != user.ID {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("unauthorized to update other's subworld template"))
		return
	}

	subworld_template.FileName = req.FileName
	subworld_template.DisplayName = req.DisplayName
	subworld_template.LevelIpfsUri = req.LevelIpfsUri
	subworld_template.LevelCentralizedUri = req.LevelCentralizedUri
	subworld_template.ThumbnailCentralizedUri = req.ThumbnailCentralizedUri
	subworld_template.Derivable = req.Derivable
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

	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authU == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var user model.User
	user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))

	var subworld_template model.SubworldTemplate
	idStr := c.Param("root_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if err := subworld_template.GetById(id); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if *subworld_template.CreatorId != user.ID {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("unauthorized to delete other's root subworld"))
		return
	}
	if err := subworld_template.Delete(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{})
}

func (ctrl *SubworldTemplateController) GetAllDeriv(c *gin.Context) {
	const (
		success = "Get Deriv Subworld Templates successfully"
		failed  = "Get Deriv Subworld Templates unsuccessfully"
	)

	var sts []model.SubworldTemplate
	var err error

	rootIdStr := c.Param("root_id")
	rootId, err := strconv.Atoi(rootIdStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	userIdStr := c.Request.URL.Query().Get("user_id")
	derivable := c.Request.URL.Query().Get("derivable")

	if userIdStr == "" {
		sts, err = model.GetAllDeriv(rootId)
		if err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}
	} else {
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}
		sts, err = model.GetDerivFromCreator(rootId, userId)
		if err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}
	}

	filtered_sts := []model.SubworldTemplate{}
	for i := range sts {
		if derivable == "" || strconv.Itoa(sts[i].Derivable) == derivable {
			filtered_sts = append(filtered_sts, sts[i])
		}
	}
	JSONReturn(c, http.StatusOK, success, gin.H{
		"subworld_templates": filtered_sts,
	})
}

func (ctrl *SubworldTemplateController) CreateDeriv(c *gin.Context) {
	const (
		success = "Create Deriv Subworld Template successfully"
		failed  = "Create Deriv Subworld Template unsuccessfully"
	)

	var req requestSchema.CreateSubworldTemplateDerive
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authU == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var user model.User
	user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))

	rootIdStr := c.Param("root_id")
	rootId, err := strconv.Atoi(rootIdStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	rootIdUInt := uint(rootId)
	var subworld_template = req.SubworldTemplate
	subworld_template.CreatorId = &user.ID
	subworld_template.ParentSubworldTemplateId = &rootIdUInt
	if err := subworld_template.Create(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"subworld_template": subworld_template,
	})
}

func (ctrl *SubworldTemplateController) UpdateDeriv(c *gin.Context) {
	const (
		success = "Update Deriv Subworld Template successfully"
		failed  = "Update Deriv Subworld Template unsuccessfully"
	)

	var req requestSchema.UpdateSubworldTemplateDerive
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authU == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var user model.User
	user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))

	var subworld_template model.SubworldTemplate
	idStr := c.Param("deriv_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if err := subworld_template.GetById(id); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if *subworld_template.CreatorId != user.ID {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("unauthorized to update other's subworld template"))
		return
	}

	subworld_template.FileName = req.FileName
	subworld_template.DisplayName = req.DisplayName
	subworld_template.ThumbnailCentralizedUri = req.ThumbnailCentralizedUri
	subworld_template.DerivativeUri = req.DerivativeUri
	subworld_template.Derivable = req.Derivable
	if err := subworld_template.Update(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"subworld_template": subworld_template,
	})
}

func (ctrl *SubworldTemplateController) DeleteDeriv(c *gin.Context) {
	const (
		success = "Delete Deriv Subworld Template successfully"
		failed  = "Delete Deriv Subworld Template unsuccessfully"
	)

	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authU == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var user model.User
	user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))

	var subworld_template model.SubworldTemplate
	idStr := c.Param("deriv_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if err := subworld_template.GetById(id); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if *subworld_template.CreatorId != user.ID {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("unauthorized to delete other's deriv subworld"))
		return
	}
	if err := subworld_template.Delete(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{})
}
