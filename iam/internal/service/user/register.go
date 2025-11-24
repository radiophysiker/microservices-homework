package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/radiophysiker/microservices-homework/iam/internal/model"
)

// Register регистрирует нового пользователя.
// Проверяет уникальность логина и email, хеширует пароль и создает пользователя в базе данных.
// Возвращает UUID созданного пользователя или ошибку.
func (s *Service) Register(ctx context.Context, info *model.UserInfo, password string) (string, error) {
	info.Email = normalizeEmail(info.Email)

	existingUser, err := s.userRepository.GetByLogin(ctx, info.Login)
	if err != nil && !errors.Is(err, model.ErrUserNotFound) {
		return "", fmt.Errorf("check login uniqueness: %w", err)
	}

	if existingUser != nil {
		return "", model.NewErrUserAlreadyExists(existingUser.UUID)
	}

	existingUser, err = s.userRepository.GetByEmail(ctx, info.Email)
	if err != nil && !errors.Is(err, model.ErrUserNotFound) {
		return "", fmt.Errorf("check email uniqueness: %w", err)
	}

	if existingUser != nil {
		return "", model.NewErrUserAlreadyExists(existingUser.UUID)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	userUUID := uuid.New().String()
	now := time.Now()

	user := &model.User{
		UUID:         userUUID,
		Info:         *info,
		PasswordHash: string(passwordHash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		switch {
		case errors.Is(err, model.ErrUserAlreadyExists):
			return "", model.NewErrUserAlreadyExists(userUUID)
		default:
			return "", fmt.Errorf("create user: %w", err)
		}
	}

	return userUUID, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
