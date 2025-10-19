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

// NewErrPartNotFound создает ошибку "деталь не найдена"
func NewErrPartNotFound(uuid string) error {
	return fmt.Errorf("%w: %s", ErrPartNotFound, uuid)
}

// NewErrInvalidUUID создает ошибку "некорректный UUID"
func NewErrInvalidUUID(uuid string) error {
	return fmt.Errorf("%w: %s", ErrInvalidUUID, uuid)
}
