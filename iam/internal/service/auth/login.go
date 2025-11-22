package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/radiophysiker/microservices-homework/iam/internal/model"
)

// Login выполняет вход пользователя.
// Проверяет логин и пароль, создает новую сессию в Redis с TTL и добавляет ее в множество сессий пользователя.
// Возвращает UUID созданной сессии или ошибку.
func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.userRepository.GetByLogin(ctx, login)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			return "", model.ErrInvalidCredentials
		default:
			return "", fmt.Errorf("get user by login: %w", err)
		}
	}

	if user == nil {
		return "", model.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", model.ErrInvalidCredentials
	}

	sessionUUID := uuid.New().String()
	now := time.Now()
	expiresAt := now.Add(s.sessionTTL)

	session := &model.Session{
		UUID:      sessionUUID,
		UserUUID:  user.UUID,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: expiresAt,
	}

	if err := s.sessionRepository.Create(ctx, session); err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	if err := s.sessionRepository.AddSessionToUserSet(ctx, user.UUID, sessionUUID); err != nil {
		return "", fmt.Errorf("add session to user set: %w", err)
	}

	return sessionUUID, nil
}
