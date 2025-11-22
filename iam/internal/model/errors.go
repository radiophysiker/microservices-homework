package model

import (
	"errors"
	"fmt"
)

var (
	// ErrUserNotFound - ошибка "пользователь не найден"
	ErrUserNotFound = errors.New("user not found")
	// ErrUserAlreadyExists - ошибка "пользователь уже существует"
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrInvalidCredentials - ошибка "неверные учетные данные"
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrSessionNotFound - ошибка "сессия не найдена"
	ErrSessionNotFound = errors.New("session not found")
	// ErrInvalidSession - ошибка "сессия недействительна"
	ErrInvalidSession = errors.New("invalid session")
)

// NewErrUserNotFound создает ошибку "пользователь не найден"
func NewErrUserNotFound(uuid string) error {
	return fmt.Errorf("%w: %s", ErrUserNotFound, uuid)
}

// NewErrUserAlreadyExists создает ошибку "пользователь уже существует"
func NewErrUserAlreadyExists(uuid string) error {
	return fmt.Errorf("%w: %s", ErrUserAlreadyExists, uuid)
}

// NewErrInvalidCredentials создает ошибку "неверные учетные данные"
func NewErrInvalidCredentials(uuid string) error {
	return fmt.Errorf("%w: %s", ErrInvalidCredentials, uuid)
}

// NewErrSessionNotFound создает ошибку "сессия не найдена"
func NewErrSessionNotFound(uuid string) error {
	return fmt.Errorf("%w: %s", ErrSessionNotFound, uuid)
}

// NewErrInvalidSession создает ошибку "сессия недействительна"
func NewErrInvalidSession(uuid string) error {
	return fmt.Errorf("%w: %s", ErrInvalidSession, uuid)
}
