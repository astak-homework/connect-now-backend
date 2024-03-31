package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	usecaseAuth "github.com/astak-homework/connect-now-backend/auth/usecase"
	"github.com/astak-homework/connect-now-backend/models"
	usecaseProfile "github.com/astak-homework/connect-now-backend/profile/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	testAccountId := "id"

	r := gin.Default()
	group := r.Group("/user")

	authUseCase := new(usecaseAuth.AuthUseCaseMock)
	profileUseCase := new(usecaseProfile.ProfileUseCaseMock)

	RegisterHTTPEndpoints(group, authUseCase, profileUseCase)

	body := `
	{
		"first_name": "testfirstname",
		"second_name": "testlastname",
		"birthdate": "2017-02-01",
		"gender": "male",
		"biography": "testbiography",
		"city": "testcity",
		"password": "testpassword"
	}
	`
	profile := &models.Profile{
		ID:        testAccountId,
		FirstName: "testfirstname",
		LastName:  "testlastname",
		BirthDate: time.Date(2017, time.February, 1, 0, 0, 0, 0, time.UTC),
		Gender:    models.GenderMale,
		Biography: "testbiography",
		City:      "testcity",
	}

	authUseCase.On("SignUp", "testpassword").Return(testAccountId, nil)
	profileUseCase.On("CreateProfile", profile).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer([]byte(body)))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestGet(t *testing.T) {
	r := gin.Default()
	group := r.Group("/user")

	authUseCase := new(usecaseAuth.AuthUseCaseMock)
	profileUseCase := new(usecaseProfile.ProfileUseCaseMock)

	RegisterHTTPEndpoints(group, authUseCase, profileUseCase)

	p := &models.Profile{
		ID:        "id",
		FirstName: "firstname",
		LastName:  "lastname",
		BirthDate: time.Date(2017, time.February, 1, 0, 0, 0, 0, time.UTC),
		Gender:    models.GenderMale,
		Biography: "biography",
		City:      "city",
	}

	profileUseCase.On("GetProfile", "id").Return(p, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/get/id", nil)
	r.ServeHTTP(w, req)

	expectedOut := `
	{
		"id": "id",
		"first_name": "firstname",
		"second_name": "lastname",
		"birthdate": "2017-02-01",
		"gender": "male",
		"biography": "biography",
		"city": "city"
	}
	`

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, expectedOut, w.Body.String())
}

func TestSearch(t *testing.T) {
	r := gin.Default()
	group := r.Group("/user")

	authUseCase := new(usecaseAuth.AuthUseCaseMock)
	profileUseCase := new(usecaseProfile.ProfileUseCaseMock)

	RegisterHTTPEndpoints(group, authUseCase, profileUseCase)

	p := &models.Profile{
		ID:        "c68c093c-ecf1-4e80-96d1-402cf9ec46cf",
		FirstName: "Люция",
		LastName:  "Меркурьева",
		BirthDate: time.Date(1982, 9, 29, 0, 0, 0, 0, time.UTC),
		Gender:    models.GnderFemale,
		Biography: "bio",
		City:      "Сусуман",
	}

	profileUseCase.On("SearchProfile", "Л", "М").Return([]*models.Profile{p}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/search?first_name=Л&last_name=М", nil)
	r.ServeHTTP(w, req)

	expectedOut := `
	[
		{
			"id": "c68c093c-ecf1-4e80-96d1-402cf9ec46cf",
			"first_name": "Люция",
			"second_name": "Меркурьева",
			"birthdate": "1982-09-29",
			"gender": "female",
			"biography": "bio",
			"city": "Сусуман"
		}
	]
	`

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, expectedOut, w.Body.String())
}
