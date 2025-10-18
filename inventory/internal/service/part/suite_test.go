package part

import (
	"context"
	"testing"

	repomocks "github.com/radiophysiker/microservices-homework/inventory/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	repo    *repomocks.MockPartRepository
	service *Service
	ctx     context.Context
}

func (s *ServiceTestSuite) SetupTest() {
	s.repo = repomocks.NewMockPartRepository(s.T())
	s.service = NewService(s.repo)
	s.ctx = context.Background()
}

func (s *ServiceTestSuite) TearDownTest() {
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
