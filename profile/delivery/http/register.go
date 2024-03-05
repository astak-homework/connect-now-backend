package http

import (
	"github.com/astak-homework/connect-now-backend/profile"
	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc profile.UseCase) {
	h := NewHandler(uc)

	profiles := router.Group("/profiles")
	{
		profiles.POST("", h.Create)
		profiles.GET("", h.Get)
		profiles.DELETE("", h.Delete)
	}
}
