package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

// RepositoryTestSuite содержит общее окружение для всех тестов репозитория
type RepositoryTestSuite struct {
	suite.Suite
	repo *Repository
	ctx  context.Context
}

// SetupTest запускается перед каждым тестом
func (s *RepositoryTestSuite) SetupTest() {
	s.repo = NewRepository()
	s.ctx = context.Background()
}

// TestRepositorySuite запускает все тесты suite
func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
