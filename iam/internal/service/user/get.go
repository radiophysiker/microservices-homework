package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/radiophysiker/microservices-homework/iam/internal/model"
)

// Get возвращает пользователя по UUID.
// Валидирует формат UUID и возвращает пользователя из базы данных или ошибку.
func (s *Service) Get(ctx context.Context, userUUID string) (*model.User, error) {
	_, err := uuid.Parse(userUUID)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid format: %w", err)
	}

	user, err := s.userRepository.GetByUUID(ctx, userUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			return nil, model.NewErrUserNotFound(userUUID)
		default:
			return nil, fmt.Errorf("get user by uuid: %w", err)
		}
	}

	if user == nil {
		return nil, model.NewErrUserNotFound(userUUID)
	}

	return user, nil
}
