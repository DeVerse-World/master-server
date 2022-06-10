package controller

import "github.com/gin-gonic/gin"

type BaseResponse struct {
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func abortWithStatusError(c *gin.Context, code int, message string, err error) {
	resp := BaseResponse{
		Message: message,
		Error:   err.Error(),
	}
	c.AbortWithStatusJSON(code, resp)
}

func JSONReturn(c *gin.Context, code int, message string, data interface{}) {
	resp := BaseResponse{
		Message: message,
		Data:    data,
	}
	c.JSON(code, resp)
}
