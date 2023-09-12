package request

import "github.com/hyperjiang/gin-skeleton/model"

type CreateSubworldTemplateRoot struct {
	SubworldTemplate model.SubworldTemplate `json:"subworld_template" binding:"required"`
}

type UpdateSubworldTemplateRoot struct {
	FileName                string `json:"file_name"`
	DisplayName             string `json:"display_name"`
	LevelIpfsUri            string `json:"level_ipfs_uri"`
	LevelCentralizedUri     string `json:"level_centralized_uri"`
	ThumbnailCentralizedUri string `json:"thumbnail_centralized_uri"`
	ImageParonamaUri        string `json:"image_paronama_uri"`
	Derivable               int    `json:"derivable"`
}

type CreateSubworldTemplateDerive struct {
	SubworldTemplate model.SubworldTemplate `json:"subworld_template" binding:"required"`
}

type UpdateSubworldTemplateDerive struct {
	FileName                string `json:"file_name"`
	DisplayName             string `json:"display_name"`
	LevelIpfsUri            string `json:"level_ipfs_uri"`
	LevelCentralizedUri     string `json:"level_centralized_uri"`
	ThumbnailCentralizedUri string `json:"thumbnail_centralized_uri"`
	DerivativeUri           string `json:"derivative_uri"`
	ImageParonamaUri        string `json:"image_paronama_uri"`
	Derivable               int    `json:"derivable"`
}

type AddSubworldTemplateTags struct {
	TagNames []string `json:"tag_names"`
}
