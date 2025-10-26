package part

import (
	"context"
	"testing"

	repositoryMocks "github.com/radiophysiker/microservices-homework/inventory/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

// RepositoryTestSuite содержит общее окружение для всех тестов репозитория
type RepositoryTestSuite struct {
	suite.Suite
	repo *repositoryMocks.MockPartRepository
	ctx  context.Context
}

// SetupTest запускается перед каждым тестом
func (s *RepositoryTestSuite) SetupTest() {
	s.repo = repositoryMocks.NewMockPartRepository(s.T())
	s.ctx = context.Background()
}

// TestRepositorySuite запускает все тесты suite
func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
