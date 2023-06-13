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

type SubworldAssociation struct {
	NftGalleries []model.NftGallery
}

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

	var a SubworldAssociation
	nftGalleries, err := model.GetAllBySubworldTemplateId(id)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	a.NftGalleries = nftGalleries
	var creator model.User
	if st.CreatorId != nil {
		creator.GetUserByIdUInt(*st.CreatorId)
	}
	JSONReturn(c, http.StatusOK, success, gin.H{
		"subworld_template": st,
		"association":       a,
		"creator_info":      creator,
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

	combined_sts, _ := ctrl.enrichSubworldTemplates(filtered_sts)

	JSONReturn(c, http.StatusOK, success, gin.H{
		"subworld_templates":          filtered_sts,
		"enriched_subworld_templates": combined_sts,
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
	subworld_template.LevelIpfsUri = req.LevelIpfsUri
	subworld_template.LevelCentralizedUri = req.LevelCentralizedUri
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

func (ctrl *SubworldTemplateController) IncrementStats(c *gin.Context) {
	const (
		success = "Increment Stats Subworld Template successfully"
		failed  = "Increment Stats Subworld Template unsuccessfully"
	)
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

	stats_type := c.Request.URL.Query().Get("type")
	if stats_type == "num_plays" {
		subworld_template.NumPlays = subworld_template.NumPlays + 1
	} else if stats_type == "num_views" {
		subworld_template.NumViews = subworld_template.NumViews + 1
	} else if stats_type == "num_clicks" {
		subworld_template.NumClicks = subworld_template.NumClicks + 1
	} else {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("invalid stats type"))
		return
	}
	if err := subworld_template.Update(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"subworld_template": subworld_template,
	})
}

func (ctrl *SubworldTemplateController) enrichSubworldTemplates(
	sts []model.SubworldTemplate,
) ([]model.EnrichedSubworldTemplate, error) {
	enriched_sts := []model.EnrichedSubworldTemplate{}
	for _, s := range sts {
		if s.CreatorId != nil {
			var creator model.User
			if err := creator.GetUserByIdUInt(*s.CreatorId); err != nil {
				return nil, err
			}
			enriched_sts = append(enriched_sts, model.EnrichedSubworldTemplate{
				Template: s,
				CreatorInfo: struct {
					Id   uint
					Name string
				}{
					Name: creator.Name,
					Id:   creator.ID,
				},
			})
		}
	}
	return enriched_sts, nil
}
