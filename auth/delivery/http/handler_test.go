package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astak-homework/connect-now-backend/auth/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSignIn(t *testing.T) {
	r := gin.Default()
	uc := new(usecase.AuthUseCaseMock)

	RegisterHTTPEndpoints(r, uc)

	signInBody := &signInput{
		AccountId: "testuser",
		Password:  "testpass",
	}
	body, err := json.Marshal(signInBody)
	assert.NoError(t, err)

	uc.On("SignIn", signInBody.AccountId, signInBody.Password).Return("jwt", nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	response := new(signInResponse)
	err = json.Unmarshal(w.Body.Bytes(), response)
	assert.NoError(t, err)
	assert.Equal(t, "jwt", response.Token)
}
