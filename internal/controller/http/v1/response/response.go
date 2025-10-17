package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}

// SendSuccess — send standardized success response
func SendSuccess(ctx *gin.Context, code int, data interface{}) {
	ctx.JSON(code, Response{
		Meta: Meta{
			Code:    code,
			Message: http.StatusText(code),
		},
		Data: data,
	})
}

// SendError — send standardized error response
func SendError(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, Response{
		Meta: Meta{
			Code:    code,
			Message: message,
		},
	})
}

// Swagger response structs for documentation

// ErrorResponse represents an error response
type ErrorResponse struct {
	Meta Meta `json:"meta"`
}

// LoginSuccessData represents the data field in login success response
type LoginSuccessData struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
}

// LoginSuccessResponse represents a successful login response
type LoginSuccessResponse struct {
	Meta Meta             `json:"meta"`
	Data LoginSuccessData `json:"data"`
}
