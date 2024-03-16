package http

import (
	"errors"
	"net/http"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	useCase auth.UseCase
}

type signInput struct {
	AccountId string `json:"id" binding:"uuid"`
	Password  string `json:"password" binding:"min=8,max=72"`
}

type signInResponse struct {
	Token string `json:"token"`
}

func NewHandler(useCase auth.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) SignIn(c *gin.Context) {
	inp := new(signInput)
	if err := c.ShouldBindJSON(inp); err != nil {
		log.Error().Err(err).Msg("auth.SignIn: couldn't bind JSON")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"description": i18n.MustGetMessage(c, "invalid_data")})
		return
	}

	token, err := h.useCase.SignIn(c.Request.Context(), inp.AccountId, inp.Password)
	if err != nil {
		log.Error().Err(err).Msg("auth.SignIn: couldn't sign in")
		if errors.Is(err, auth.ErrUserNotFound) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, signInResponse{Token: token})
}
