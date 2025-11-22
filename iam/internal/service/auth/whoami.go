package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/radiophysiker/microservices-homework/iam/internal/model"
)

// Whoami возвращает информацию о текущей сессии и пользователе.
// Проверяет существование сессии, ее срок действия и статус отзыва.
// Возвращает сессию и пользователя или ошибку.
func (s *Service) Whoami(ctx context.Context, sessionUUID string) (*model.Session, *model.User, error) {
	session, err := s.sessionRepository.Get(ctx, sessionUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidCredentials):
			return nil, nil, model.NewErrSessionNotFound(sessionUUID)
		default:
			return nil, nil, fmt.Errorf("get session: %w", err)
		}
	}

	if session == nil {
		return nil, nil, model.ErrSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, nil, model.ErrInvalidSession
	}

	// Проверка, что сессия не отозвана
	if session.RevokedAt != nil && !session.RevokedAt.IsZero() {
		return nil, nil, model.ErrInvalidSession
	}

	user, err := s.userService.Get(ctx, session.UserUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			return nil, nil, fmt.Errorf("user not found for session: %w", err)
		default:
			return nil, nil, fmt.Errorf("get user: %w", err)
		}
	}

	if user == nil {
		return nil, nil, fmt.Errorf("user is nil for session")
	}

	return session, user, nil
}
