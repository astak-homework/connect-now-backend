package auth

import "errors"

var (
	ErrUserNotFound              = errors.New("user not found")
	ErrInvalidAccessToken        = errors.New("invalid access token")
	ErrMismatchedHashAndPassword = errors.New("mismatched hash and password")
)
