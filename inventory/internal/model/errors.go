package model

import (
	"errors"
	"fmt"
)

var (
	// ErrPartNotFound - ошибка "деталь не найдена"
	ErrPartNotFound = errors.New("part not found")
	// ErrInvalidUUID - ошибка "некорректный UUID"
	ErrInvalidUUID = errors.New("invalid UUID")
)

// NewPartNotFoundError создает ошибку "деталь не найдена"
func NewPartNotFoundError(uuid string) error {
	return fmt.Errorf("%w: %s", ErrPartNotFound, uuid)
}

// NewInvalidUUIDError создает ошибку "некорректный UUID"
func NewInvalidUUIDError(uuid string) error {
	return fmt.Errorf("%w: %s", ErrInvalidUUID, uuid)
}
