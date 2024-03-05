package http

import (
	"errors"
	"net/http"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase auth.UseCase
}

type signInput struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token"`
}

func NewHandler(useCase auth.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) SignUp(c *gin.Context) {
	inp := new(signInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := h.useCase.SignUp(c.Request.Context(), inp.UserName, inp.Password); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) SignIn(c *gin.Context) {
	inp := new(signInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := h.useCase.SignIn(c.Request.Context(), inp.UserName, inp.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, signInResponse{Token: token})
}
