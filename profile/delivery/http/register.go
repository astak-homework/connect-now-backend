package http

import (
	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/astak-homework/connect-now-backend/profile"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, authUseCase auth.UseCase, profileUseCase profile.UseCase) {
	h := NewHandler(authUseCase, profileUseCase)

	router.POST("/register", h.Create)
	router.GET("/get/:id", h.Get)
	router.GET("/search", h.Search)
}
