package http

import (
	"net/http"
	"time"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/astak-homework/connect-now-backend/errors"
	"github.com/astak-homework/connect-now-backend/models"
	"github.com/astak-homework/connect-now-backend/profile"
	"github.com/gin-gonic/gin"
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

type searchInput struct {
	FirstName string `form:"first_name" binding:"required"`
	LastName  string `form:"last_name" binding:"required"`
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
		c.Error(errors.NewBadRequest().Err(err).Log("profile.Create: couldn't bind JSON").InvalidData())
		return
	}

	accountId, err := h.authUseCase.SignUp(c.Request.Context(), inp.Password)
	if err != nil {
		c.Error(errors.NewInternal().Err(err).Log("profile.Create: couldn't sign up").Internal())
		return
	}

	model, err := toModel(accountId, inp)
	if err != nil {
		c.Error(errors.NewInternal().Err(err).Log("profile.Create: couldn't convert input to model").Internal())
		return
	}

	if err := h.profileUseCase.CreateProfile(c.Request.Context(), model); err != nil {
		c.Error(errors.NewInternal().Err(err).Log("profile.Create: couldn't create profile").Internal())
		return
	}

	c.IndentedJSON(http.StatusOK, &createResponse{AccountId: accountId})
}

func (h *Handler) Get(c *gin.Context) {
	accountId := c.Param("id")

	profile, err := h.profileUseCase.GetProfile(c.Request.Context(), accountId)
	if err != nil {
		c.Error(errors.NewInternal().Err(err).Log("profile.Get: couldn't get profile").Internal())
		return
	}

	c.IndentedJSON(http.StatusOK, toGetResponse(profile))
}

func (h *Handler) Delete(c *gin.Context) {
	inp := new(deleteInput)
	if err := c.ShouldBindJSON(inp); err != nil {
		c.Error(errors.NewBadRequest().Err(err).Log("profile.Delete: couldn't bind JSON").InvalidData())
		return
	}

	if err := h.profileUseCase.DeleteProfile(c.Request.Context(), inp.ID); err != nil {
		c.Error(errors.NewInternal().Err(err).Log("profile.Delete: couldn't delete profile").Internal())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) Search(c *gin.Context) {
	inp := new(searchInput)

	if err := c.ShouldBind(&inp); err != nil {
		c.Error(errors.NewBadRequest().Err(err).Log("profile.Search: couldn't bind query string").InvalidData())
		return
	}

	profiles, err := h.profileUseCase.SearchProfile(c.Request.Context(), inp.FirstName, inp.LastName)
	if err != nil {
		c.Error(errors.NewInternal().Err(err).Log("profile.Search: couldn't search profiles").Internal())
		return
	}

	c.IndentedJSON(http.StatusOK, toSearchResponse(profiles))
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

func toSearchResponse(ps []*models.Profile) []*getResponse {
	response := []*getResponse{}
	for _, p := range ps {
		response = append(response, toGetResponse(p))
	}
	return response
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
