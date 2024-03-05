package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/astak-homework/connect-now-backend/models"
	"github.com/astak-homework/connect-now-backend/profile/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	testLogin := &models.Login{
		ID:       "id",
		UserName: "testUser",
		Password: "testpass",
	}

	r := gin.Default()
	group := r.Group("/api", func(c *gin.Context) {
		c.Set(auth.CtxLoginKey, testLogin)
	})

	uc := new(usecase.ProfileUseCaseMock)

	RegisterHTTPEndpoints(group, uc)

	inp := &createInput{
		FirstName: "testfirstname",
		LastName:  "testlastname",
		BirthDate: time.Now(),
		Gender:    models.GenderMale,
		Biography: "testbiography",
		City:      "testcity",
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	birthDateMock := mock.MatchedBy(func(birthDate time.Time) bool { return birthDate.Equal(inp.BirthDate) })
	uc.On("CreateProfile", testLogin.ID, inp.FirstName, inp.LastName, birthDateMock, inp.Gender, inp.Biography, inp.City).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/profiles", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestGet(t *testing.T) {
	accountId := "id"
	birthDate := time.Now()

	testAccount := &models.Login{
		ID:       accountId,
		UserName: "testuser",
	}

	r := gin.Default()
	group := r.Group("/api", func(c *gin.Context) {
		c.Set(auth.CtxLoginKey, testAccount)
	})

	uc := new(usecase.ProfileUseCaseMock)

	RegisterHTTPEndpoints(group, uc)

	p := &models.Profile{
		ID:        accountId,
		FirstName: "firstname",
		LastName:  "lastname",
		BirthDate: birthDate,
		Gender:    models.GenderMale,
		Biography: "biography",
		City:      "city",
	}

	uc.On("GetProfile", accountId).Return(p, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/profiles", nil)
	r.ServeHTTP(w, req)

	expectedOut := &getResponse{
		ID:        accountId,
		FirstName: "firstname",
		LastName:  "lastname",
		BirthDate: birthDate,
		Gender:    models.GenderMale,
		Biography: "biography",
		City:      "city",
	}

	expectedOutBody, err := json.Marshal(expectedOut)
	assert.NoError(t, err)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, string(expectedOutBody), w.Body.String())
}

func TestDelete(t *testing.T) {
	accountId := "id"

	testAccount := &models.Login{
		ID:       accountId,
		UserName: "testuser",
	}

	r := gin.Default()
	group := r.Group("/api", func(c *gin.Context) {
		c.Set(auth.CtxLoginKey, testAccount)
	})

	uc := new(usecase.ProfileUseCaseMock)

	RegisterHTTPEndpoints(group, uc)

	inp := &deleteInput{
		ID: accountId,
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc.On("DeleteProfile", accountId).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/profiles", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
