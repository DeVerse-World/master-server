package request

import "github.com/hyperjiang/gin-skeleton/model"

type CreateAvatar struct {
	Avatar model.Avatar `json:"avatar" binding:"required"`
}

type UpdateAvatar struct {
	PostprocessUrl string `json:"postprocess_url"`
	Name           string `json:"name"`
}
