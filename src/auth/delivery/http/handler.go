package http

import (
	"net/http"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/astak-homework/connect-now-backend/errors"
	"github.com/gin-gonic/gin"
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
		c.Error(errors.NewBadRequest().Err(err).Log("auth.SignIn: couldn't bind JSON").InvalidData())
		return
	}

	token, err := h.useCase.SignIn(c.Request.Context(), inp.AccountId, inp.Password)
	if err != nil {
		logMsg := "auth.SignIn: couldn't sign in"
		if errors.Is(err, auth.ErrUserNotFound) {
			c.Error(errors.NewNotFound().Err(err).Log(logMsg).UserNotFound())
			return
		}
		if errors.Is(err, auth.ErrMismatchedHashAndPassword) {
			c.Error(errors.NewBadRequest().Err(err).Log(logMsg).IncorrectPassword())
			return
		}
		c.Error(errors.NewInternal().Err(err).Log(logMsg).Internal())
		return
	}

	c.IndentedJSON(http.StatusOK, signInResponse{Token: token})
}
