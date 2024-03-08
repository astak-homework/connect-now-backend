package usecase

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/astak-homework/connect-now-backend/config"
	"github.com/astak-homework/connect-now-backend/models"
	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	jwt.RegisteredClaims
	AccountID string `json:"account_id"`
	UserName  string `json:"user_name"`
}

type AuthUseCase struct {
	loginRepo      auth.LoginRepository
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthUseCase(loginRepo auth.LoginRepository, cfg *config.Auth) *AuthUseCase {
	return &AuthUseCase{
		loginRepo:      loginRepo,
		hashSalt:       cfg.HashSalt,
		signingKey:     []byte(cfg.SigningKey),
		expireDuration: time.Second * cfg.TokenTTL,
	}
}

func (a *AuthUseCase) SignUp(ctx context.Context, username, password string) error {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))

	login := &models.Login{
		UserName: username,
		Password: fmt.Sprintf("%x", pwd.Sum(nil)),
	}

	return a.loginRepo.CreateLogin(ctx, login)
}

func (a *AuthUseCase) SignIn(ctx context.Context, username, password string) (string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))
	password = fmt.Sprintf("%x", pwd.Sum(nil))

	login, err := a.loginRepo.GetLogin(ctx, username, password)
	if err != nil {
		return "", err
	}

	claims := AuthClaims{
		AccountID: login.ID,
		UserName:  login.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.expireDuration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.signingKey)
}

func (a *AuthUseCase) ParseToken(ctx context.Context, accessToken string) (*models.Login, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpectged signing method: %v", token.Header["alg"])
		}
		return a.signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return &models.Login{
			ID:       claims.AccountID,
			UserName: claims.UserName,
		}, nil
	}

	return nil, auth.ErrInvalidAccessToken
}
