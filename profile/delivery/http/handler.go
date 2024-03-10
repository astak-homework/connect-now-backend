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
	FirstName string           `json:"first_name"`
	LastName  string           `json:"second_name"`
	BirthDate models.CivilTime `json:"birthdate"`
	Gender    models.Gender    `json:"gender"`
	Biography string           `json:"biography"`
	City      string           `json:"city"`
	Password  string           `json:"password"`
}

type createResponse struct {
	AccountId string `json:"user_id"`
}

type getResponse struct {
	ID        string           `json:"id"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"second_name"`
	BirthDate models.CivilTime `json:"birthdate"`
	Gender    models.Gender    `json:"gender"`
	Biography string           `json:"biography"`
	City      string           `json:"city"`
}

type deleteInput struct {
	ID string `json:"id"`
}

type Handler struct {
	authUseCase    auth.UseCase
	profileUseCase profile.UseCase
}

func NewHandler(authUseCase auth.UseCase, profileUseCase profile.UseCase) *Handler {
	return &Handler{
		authUseCase:    authUseCase,
		profileUseCase: profileUseCase,
	}
}

func (h *Handler) Create(c *gin.Context) {
	inp := new(createInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	accountId, err := h.authUseCase.SignUp(c.Request.Context(), inp.Password)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	if err := h.profileUseCase.CreateProfile(c.Request.Context(), toModel(accountId, inp)); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, &createResponse{AccountId: accountId})
}

func (h *Handler) Get(c *gin.Context) {
	accountId := c.Param("id")

	profile, err := h.profileUseCase.GetProfile(c.Request.Context(), accountId)
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

	if err := h.profileUseCase.DeleteProfile(c.Request.Context(), inp.ID); err != nil {
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
		BirthDate: models.CivilTime(p.BirthDate),
		Gender:    p.Gender,
		Biography: p.Biography,
		City:      p.City,
	}
}

func toModel(accountId string, i *createInput) *models.Profile {
	return &models.Profile{
		ID:        accountId,
		FirstName: i.FirstName,
		LastName:  i.LastName,
		BirthDate: time.Time(i.BirthDate),
		Gender:    i.Gender,
		Biography: i.Biography,
		City:      i.City,
	}
}
