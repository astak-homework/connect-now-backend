package http

import (
	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, uc auth.UseCase) {
	h := NewHandler(uc)
	router.POST("/login", h.SignIn)
}
