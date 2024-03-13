package http

import (
	"net/http"
	"time"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/astak-homework/connect-now-backend/models"
	"github.com/astak-homework/connect-now-backend/profile"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Msg("profile.Create: couldn't bind JSON")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	accountId, err := h.authUseCase.SignUp(c.Request.Context(), inp.Password)
	if err != nil {
		log.Error().Err(err).Msg("profile.Create: couldn't sing up")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := h.profileUseCase.CreateProfile(c.Request.Context(), toModel(accountId, inp)); err != nil {
		log.Error().Err(err).Msg("profile.Create: couldn't create profile")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, &createResponse{AccountId: accountId})
}

func (h *Handler) Get(c *gin.Context) {
	accountId := c.Param("id")

	profile, err := h.profileUseCase.GetProfile(c.Request.Context(), accountId)
	if err != nil {
		log.Error().Err(err).Msg("profile.Get: couldn't get profile")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, toGetResponse(profile))
}

func (h *Handler) Delete(c *gin.Context) {
	inp := new(deleteInput)
	if err := c.BindJSON(inp); err != nil {
		log.Error().Err(err).Msg("profile.Delete: couldn't bind JSON")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := h.profileUseCase.DeleteProfile(c.Request.Context(), inp.ID); err != nil {
		log.Error().Err(err).Msg("profile.Delete: couldn't delete profile")
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
