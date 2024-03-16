package http

import (
	"net/http"
	"time"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/astak-homework/connect-now-backend/models"
	"github.com/astak-homework/connect-now-backend/profile"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type createInput struct {
	FirstName string        `json:"first_name" binding:"required,max=50"`
	LastName  string        `json:"second_name" binding:"required,max=50"`
	BirthDate string        `json:"birthdate" binding:"datetime=2006-01-02"`
	Gender    models.Gender `json:"gender" binding:"oneof=male female"`
	Biography string        `json:"biography" binding:"required"`
	City      string        `json:"city" binding:"required,max=50"`
	Password  string        `json:"password" binding:"required,max=72"`
}

type createResponse struct {
	AccountId string `json:"user_id"`
}

type getResponse struct {
	ID        string        `json:"id"`
	FirstName string        `json:"first_name"`
	LastName  string        `json:"second_name"`
	BirthDate string        `json:"birthdate"`
	Gender    models.Gender `json:"gender"`
	Biography string        `json:"biography"`
	City      string        `json:"city"`
}

type deleteInput struct {
	ID string `json:"id" binding:"uuid"`
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
	if err := c.ShouldBindJSON(inp); err != nil {
		log.Error().Err(err).Msg("profile.Create: couldn't bind JSON")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"description": i18n.MustGetMessage(c, "invalid_data")})
		return
	}

	accountId, err := h.authUseCase.SignUp(c.Request.Context(), inp.Password)
	if err != nil {
		log.Error().Err(err).Msg("profile.Create: couldn't sing up")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	model, err := toModel(accountId, inp)
	if err != nil {
		log.Error().Err(err).Msg("profile.Create: couldn't convert input to model")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := h.profileUseCase.CreateProfile(c.Request.Context(), model); err != nil {
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
	if err := c.ShouldBindJSON(inp); err != nil {
		log.Error().Err(err).Msg("profile.Delete: couldn't bind JSON")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"description": i18n.MustGetMessage(c, "invalid_data")})
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
		BirthDate: time.Time(p.BirthDate).Format("2006-01-02"),
		Gender:    p.Gender,
		Biography: p.Biography,
		City:      p.City,
	}
}

func toModel(accountId string, i *createInput) (*models.Profile, error) {
	birthDate, err := time.Parse("2006-01-02", i.BirthDate)
	if err != nil {
		return nil, err
	}
	return &models.Profile{
		ID:        accountId,
		FirstName: i.FirstName,
		LastName:  i.LastName,
		BirthDate: birthDate,
		Gender:    i.Gender,
		Biography: i.Biography,
		City:      i.City,
	}, nil
}
