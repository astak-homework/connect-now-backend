package http

import (
	"net/http"
	"time"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/astak-homework/connect-now-backend/models"
	"github.com/astak-homework/connect-now-backend/profile"
	"github.com/gin-gonic/gin"
)

type createInput struct {
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
	BirthDate time.Time     `json:"birth_date"`
	Gender    models.Gender `json:"gender"`
	Biography string        `json:"biography"`
	City      string        `json:"city"`
}

type getResponse struct {
	ID        string        `jsin:"id"`
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
	BirthDate time.Time     `json:"birth_date"`
	Gender    models.Gender `json:"gender"`
	Biography string        `json:"biography"`
	City      string        `json:"city"`
}

type deleteInput struct {
	ID string `json:"id"`
}

type Handler struct {
	useCase profile.UseCase
}

func NewHandler(useCase profile.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) Create(c *gin.Context) {
	inp := new(createInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	login := c.MustGet(auth.CtxLoginKey).(*models.Login)

	if err := h.useCase.CreateProfile(c.Request.Context(), login.ID, inp.FirstName, inp.LastName, inp.BirthDate, inp.Gender, inp.Biography, inp.City); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) Get(c *gin.Context) {
	login := c.MustGet(auth.CtxLoginKey).(*models.Login)

	profile, err := h.useCase.GetProfile(c.Request.Context(), login.ID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, toGetResponse(profile))
}

func (h *Handler) Delete(c *gin.Context) {
	inp := new(deleteInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := h.useCase.DeleteProfile(c.Request.Context(), inp.ID); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func toGetResponse(p *models.Profile) *getResponse {
	return &getResponse{
		ID:        p.ID,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		BirthDate: p.BirthDate,
		Gender:    p.Gender,
		Biography: p.Biography,
		City:      p.City,
	}
}
