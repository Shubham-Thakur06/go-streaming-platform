package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccessDenied       = errors.New("access denied")
	ErrMediaNotFound      = errors.New("media not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrUnauthorized       = errors.New("unauthorized")
)
