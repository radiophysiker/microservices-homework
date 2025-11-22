package auth

import (
	"time"

	"github.com/radiophysiker/microservices-homework/iam/internal/repository"
	"github.com/radiophysiker/microservices-homework/iam/internal/service"
)

// Service реализует интерфейс AuthService
type Service struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	userService       service.UserService
	sessionTTL        time.Duration
}

// NewService создает новый экземпляр Service
func NewService(
	userRepository repository.UserRepository,
	sessionRepository repository.SessionRepository,
	userService service.UserService,
	sessionTTL time.Duration,
) *Service {
	return &Service{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		userService:       userService,
		sessionTTL:        sessionTTL,
	}
}
