package order

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
)

// RepositoryTestSuite содержит общее окружение для всех тестов репозитория
type RepositoryTestSuite struct {
	suite.Suite
	repo *Repository
	ctx  context.Context

	testOrderUUID uuid.UUID
	testUserUUID  uuid.UUID
	testPartUUID  uuid.UUID
}

// SetupTest запускается перед каждым тестом
func (s *RepositoryTestSuite) SetupTest() {
	s.repo = NewRepository()
	s.ctx = context.Background()

	s.testOrderUUID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440003")
	s.testUserUUID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	s.testPartUUID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")
}

// createTestOrder создает тестовый заказ с заданными параметрами
func (s *RepositoryTestSuite) createTestOrder(orderUUID uuid.UUID, totalPrice float64, status model.Status) *model.Order {
	return &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   s.testUserUUID,
		PartUUIDs:  []uuid.UUID{s.testPartUUID},
		TotalPrice: totalPrice,
		Status:     status,
	}
}

// createTestOrderWithParts создает тестовый заказ с множественными частями
func (s *RepositoryTestSuite) createTestOrderWithParts(orderUUID uuid.UUID, partUUIDs []uuid.UUID, totalPrice float64, status model.Status) *model.Order {
	return &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   s.testUserUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     status,
	}
}

// createTestOrderWithPayment создает тестовый заказ с информацией о платеже
func (s *RepositoryTestSuite) createTestOrderWithPayment(orderUUID uuid.UUID, totalPrice float64, status model.Status, transactionUUID *uuid.UUID, paymentMethod *model.PaymentMethod) *model.Order {
	return &model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        s.testUserUUID,
		PartUUIDs:       []uuid.UUID{s.testPartUUID},
		TotalPrice:      totalPrice,
		Status:          status,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
	}
}

// ptrUUID - хелпер для создания указателя на UUID
func ptrUUID(u uuid.UUID) *uuid.UUID {
	return &u
}

// ptrPaymentMethod - хелпер для создания указателя на PaymentMethod
func ptrPaymentMethod(pm model.PaymentMethod) *model.PaymentMethod {
	return &pm
}

// TestRepositorySuite запускает все тесты suite
func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
