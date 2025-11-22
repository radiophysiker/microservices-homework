package user

import (
	"github.com/radiophysiker/microservices-homework/iam/internal/repository"
)

// Service реализует интерфейс UserService
type Service struct {
	userRepository repository.UserRepository
}

// NewService создает новый экземпляр Service
func NewService(userRepository repository.UserRepository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}
