package response

import (
	"github.com/gin-gonic/gin"
)

type APIError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type APISuccess struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omiempty"`
}

func Error(c *gin.Context, statuscode int, message string) {
	c.JSON(statuscode, APIError{
		Status:  "error",
		Message: message,
		Code:    statuscode,
	})
}

func Success(c *gin.Context, statuscode int, message string, data interface{}) {
	c.JSON(statuscode, APISuccess{
		Status:  "success",
		Message: message,
		Code:    statuscode,
		Data:    data,
	})
}
