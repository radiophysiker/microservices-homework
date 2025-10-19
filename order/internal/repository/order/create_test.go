package order

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	orderv1 "github.com/radiophysiker/microservices-homework/shared/pkg/openapi/order/v1"
)

// TestCreateOrder проверяет создание заказа
func (s *RepositoryTestSuite) TestCreateOrder() {
	tests := []struct {
		name      string
		order     *model.Order
		setupFunc func()
	}{
		{
			name:      "success",
			order:     s.createTestOrder(s.testOrderUUID, 150.75, orderv1.OrderStatusPENDINGPAYMENT),
			setupFunc: nil,
		},
		{
			name: "success_with_multiple_parts",
			order: s.createTestOrderWithParts(
				uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"),
				[]uuid.UUID{
					s.testPartUUID,
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"),
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"),
				},
				500.00,
				orderv1.OrderStatusPENDINGPAYMENT,
			),
			setupFunc: nil,
		},
		{
			name:  "success_overwrite_existing",
			order: s.createTestOrder(s.testOrderUUID, 200.00, orderv1.OrderStatusPAID),
			setupFunc: func() {
				existingOrder := s.createTestOrder(s.testOrderUUID, 100.00, orderv1.OrderStatusPENDINGPAYMENT)
				err := s.repo.CreateOrder(s.ctx, existingOrder)
				require.NoError(s.T(), err)
			},
		},
		{
			name: "success_with_payment_info",
			order: s.createTestOrderWithPayment(
				uuid.MustParse("550e8400-e29b-41d4-a716-446655440007"),
				300.00,
				orderv1.OrderStatusPENDINGPAYMENT,
				nil,
				ptrPaymentMethod(orderv1.OrderDtoPaymentMethodCARD),
			),
			setupFunc: nil,
		},
		{
			name: "empty_parts",
			order: &model.Order{
				OrderUUID:  uuid.MustParse("550e8400-e29b-41d4-a716-446655440008"),
				UserUUID:   s.testUserUUID,
				PartUUIDs:  []uuid.UUID{}, // Пустой список частей
				TotalPrice: 100.00,
				Status:     orderv1.OrderStatusPENDINGPAYMENT,
			},
			setupFunc: nil,
		},
		{
			name: "negative_price",
			order: &model.Order{
				OrderUUID:  uuid.MustParse("550e8400-e29b-41d4-a716-446655440009"),
				UserUUID:   s.testUserUUID,
				PartUUIDs:  []uuid.UUID{s.testPartUUID},
				TotalPrice: -50.00, // Отрицательная цена
				Status:     orderv1.OrderStatusPENDINGPAYMENT,
			},
			setupFunc: nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.setupFunc != nil {
				tt.setupFunc()
			}

			err := s.repo.CreateOrder(s.ctx, tt.order)

			require.NoError(s.T(), err)

			createdOrder, getErr := s.repo.GetOrder(s.ctx, tt.order.OrderUUID.String())
			require.NoError(s.T(), getErr)
			require.NotNil(s.T(), createdOrder)
			require.Equal(s.T(), tt.order.OrderUUID, createdOrder.OrderUUID)
			require.Equal(s.T(), tt.order.UserUUID, createdOrder.UserUUID)
			require.Equal(s.T(), tt.order.TotalPrice, createdOrder.TotalPrice)
			require.Equal(s.T(), tt.order.Status, createdOrder.Status)
			require.Equal(s.T(), len(tt.order.PartUUIDs), len(createdOrder.PartUUIDs))
		})
	}
}
