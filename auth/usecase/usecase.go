package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/astak-homework/connect-now-backend/config"
	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	jwt.RegisteredClaims
	AccountID string `json:"account_id"`
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

func (a *AuthUseCase) SignUp(ctx context.Context, password string) (string, error) {
	pwd := sha256.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))

	passwordHash := fmt.Sprintf("%x", pwd.Sum(nil))

	return a.loginRepo.CreateLogin(ctx, passwordHash)
}

func (a *AuthUseCase) SignIn(ctx context.Context, accountId, password string) (string, error) {
	pwd := sha256.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))
	password_hash := fmt.Sprintf("%x", pwd.Sum(nil))

	err := a.loginRepo.AuthenticateLogin(ctx, accountId, password_hash)
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
