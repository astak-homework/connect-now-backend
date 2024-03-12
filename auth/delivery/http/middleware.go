package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AuthMiddleware struct {
	usecase auth.UseCase
}

func NewAuthMiddleware(usecase auth.UseCase) gin.HandlerFunc {
	return (&AuthMiddleware{
		usecase: usecase,
	}).Handle
}

func (m *AuthMiddleware) Handle(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		log.Error().Msg("authMiddleware: couldn't get header")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		log.Error().Msg("authMiddleware: couldn't split auth header")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if headerParts[0] != "Bearer" {
		log.Error().Msg("authMiddleware: invalid auth header")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	login, err := m.usecase.ParseToken(c.Request.Context(), headerParts[1])
	if err != nil {
		log.Error().Err(err).Msg("authMiddleware: could't parse token")
		status := http.StatusInternalServerError
		if errors.Is(err, auth.ErrInvalidAccessToken) {
			status = http.StatusUnauthorized
		}

		c.AbortWithStatus(status)
		return
	}

	c.Set(auth.CtxLoginKey, login)
}
