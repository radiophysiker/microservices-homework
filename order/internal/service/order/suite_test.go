package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	clientmocks "github.com/radiophysiker/microservices-homework/order/internal/client/grpc/mocks"
	repomocks "github.com/radiophysiker/microservices-homework/order/internal/repository/mocks"
)

// ServiceTestSuite содержит общее окружение для всех тестов сервиса
type ServiceTestSuite struct {
	suite.Suite
	repo            *repomocks.MockOrderRepository
	inventoryClient *clientmocks.MockInventoryClient
	paymentClient   *clientmocks.MockPaymentClient
	service         *Service
	ctx             context.Context
}

// SetupTest запускается перед каждым тестом
func (s *ServiceTestSuite) SetupTest() {
	s.repo = repomocks.NewMockOrderRepository(s.T())
	s.inventoryClient = clientmocks.NewMockInventoryClient(s.T())
	s.paymentClient = clientmocks.NewMockPaymentClient(s.T())
	s.service = NewService(s.repo, s.inventoryClient, s.paymentClient, nil)
	s.ctx = context.Background()
}

// TestServiceSuite запускает все тесты suite
func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
