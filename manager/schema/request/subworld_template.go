package request

import "github.com/hyperjiang/gin-skeleton/model"

type CreateSubworldTemplateRoot struct {
	SubworldTemplate model.SubworldTemplate `json:"subworld_template" binding:"required"`
}

type UpdateSubworldTemplateRoot struct {
	FileName                string `json:"file_name"  binding:"required"`
	DisplayName             string `json:"display_name"  binding:"required"`
	LevelIpfsUri            string `json:"level_ipfs_uri"  binding:"required"`
	LevelCentralizedUri     string `json:"level_centralized_uri"  binding:"required"`
	ThumbnailCentralizedUri string `json:"thumbnail_centralized_uri"  binding:"required"`
}

type CreateSubworldTemplateDerive struct {
	SubworldTemplate model.SubworldTemplate `json:"subworld_template" binding:"required"`
}

type UpdateSubworldTemplateDerive struct {
	FileName                string `json:"file_name"`
	DisplayName             string `json:"display_name"`
	ThumbnailCentralizedUri string `json:"thumbnail_centralized_uri"`
	DerivativeUri           string `json:"derivative_uri"`
}
