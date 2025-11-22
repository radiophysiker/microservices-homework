package repository

import (
	"context"

	"github.com/radiophysiker/microservices-homework/iam/internal/model"
)

// UserRepository описывает операции с пользователями в PostgreSQL.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByUUID(ctx context.Context, uuid string) (*model.User, error)
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

// SessionRepository описывает операции с сессиями в Redis.
type SessionRepository interface {
	Create(ctx context.Context, session *model.Session) error
	Get(ctx context.Context, sessionUUID string) (*model.Session, error)
	AddSessionToUserSet(ctx context.Context, userUUID, sessionUUID string) error
}
