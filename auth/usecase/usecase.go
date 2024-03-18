package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/astak-homework/connect-now-backend/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthClaims struct {
	jwt.RegisteredClaims
	AccountID string `json:"account_id"`
}

type AuthUseCase struct {
	loginRepo      auth.LoginRepository
	hashCost       int
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthUseCase(loginRepo auth.LoginRepository, cfg *config.Auth) *AuthUseCase {
	return &AuthUseCase{
		loginRepo:      loginRepo,
		hashCost:       cfg.HashCost,
		signingKey:     []byte(cfg.SigningKey),
		expireDuration: time.Second * cfg.TokenTTL,
	}
}

func (a *AuthUseCase) SignUp(ctx context.Context, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), a.hashCost)
	if err != nil {
		return "", err
	}

	return a.loginRepo.CreateLogin(ctx, string(bytes))
}

func (a *AuthUseCase) SignIn(ctx context.Context, accountId, password string) (string, error) {
	passwordHash, err := a.loginRepo.GetPasswordHash(ctx, accountId)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", auth.ErrMismatchedHashAndPassword
	}
	if err != nil {
		return "", err
	}

	claims := AuthClaims{
		AccountID: accountId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.expireDuration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.signingKey)
}

func (a *AuthUseCase) ParseToken(ctx context.Context, accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpectged signing method: %v", token.Header["alg"])
		}
		return a.signingKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims.AccountID, nil
	}

	return "", auth.ErrInvalidAccessToken
}
