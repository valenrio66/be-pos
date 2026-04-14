package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valenrio66/be-pos/internal/delivery/http/dto"
)

func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, dto.APIResponse{
		Code:    code,
		Status:  http.StatusText(code),
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string, err interface{}) {
	c.JSON(code, dto.APIResponse{
		Code:    code,
		Status:  http.StatusText(code),
		Message: message,
		Errors:  err,
	})
}
