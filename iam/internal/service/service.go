package service

import (
	"context"

	"github.com/radiophysiker/microservices-homework/iam/internal/model"
)

// UserService представляет интерфейс для работы с пользователями
type UserService interface {
	// Register регистрирует нового пользователя
	Register(ctx context.Context, info *model.UserInfo, password string) (string, error)
	// Get возвращает пользователя по UUID
	Get(ctx context.Context, uuid string) (*model.User, error)
}

// AuthService представляет интерфейс для аутентификации и авторизации
type AuthService interface {
	// Login выполняет вход пользователя
	Login(ctx context.Context, login, password string) (string, error)
	// Whoami возвращает информацию о текущей сессии и пользователе
	Whoami(ctx context.Context, sessionUUID string) (*model.Session, *model.User, error)
}
