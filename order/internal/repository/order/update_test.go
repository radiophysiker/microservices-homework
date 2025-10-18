package order

import (
	"github.com/google/uuid"
	"github.com/radiophysiker/microservices-homework/order/internal/model"
	orderv1 "github.com/radiophysiker/microservices-homework/shared/pkg/openapi/order/v1"
	"github.com/stretchr/testify/require"
)

// TestUpdateOrder проверяет обновление заказа
func (s *RepositoryTestSuite) TestUpdateOrder() {
	tests := []struct {
		name      string
		order     *model.Order
		setupFunc func()
		errType   error
	}{
		{
			name:  "success",
			order: s.createTestOrder(s.testOrderUUID, 250.00, orderv1.OrderStatusPAID),
			setupFunc: func() {
				existingOrder := s.createTestOrder(s.testOrderUUID, 100.00, orderv1.OrderStatusPENDINGPAYMENT)
				err := s.repo.CreateOrder(s.ctx, existingOrder)
				require.NoError(s.T(), err)
			},
			errType: nil,
		},
		{
			name:  "success_change_status",
			order: s.createTestOrder(s.testOrderUUID, 100.00, orderv1.OrderStatusCANCELLED),
			setupFunc: func() {
				existingOrder := s.createTestOrder(s.testOrderUUID, 100.00, orderv1.OrderStatusPENDINGPAYMENT)
				err := s.repo.CreateOrder(s.ctx, existingOrder)
				require.NoError(s.T(), err)
			},
			errType: nil,
		},
		{
			name: "success_add_parts",
			order: s.createTestOrderWithParts(
				s.testOrderUUID,
				[]uuid.UUID{
					s.testPartUUID,
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"),
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"),
				},
				500.00,
				orderv1.OrderStatusPENDINGPAYMENT,
			),
			setupFunc: func() {
				existingOrder := s.createTestOrder(s.testOrderUUID, 100.00, orderv1.OrderStatusPENDINGPAYMENT)
				err := s.repo.CreateOrder(s.ctx, existingOrder)
				require.NoError(s.T(), err)
			},
			errType: nil,
		},
		{
			name: "success_add_payment_info",
			order: s.createTestOrderWithPayment(
				s.testOrderUUID,
				100.00,
				orderv1.OrderStatusPAID,
				ptrUUID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440008")),
				ptrPaymentMethod(orderv1.OrderDtoPaymentMethodCARD),
			),
			setupFunc: func() {
				existingOrder := s.createTestOrder(s.testOrderUUID, 100.00, orderv1.OrderStatusPENDINGPAYMENT)
				err := s.repo.CreateOrder(s.ctx, existingOrder)
				require.NoError(s.T(), err)
			},
			errType: nil,
		},
		{
			name:      "not_found",
			order:     s.createTestOrder(uuid.MustParse("550e8400-e29b-41d4-a716-999999999999"), 100.00, orderv1.OrderStatusPENDINGPAYMENT),
			setupFunc: nil,
			errType:   model.ErrOrderNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.setupFunc != nil {
				tt.setupFunc()
			}

			err := s.repo.UpdateOrder(s.ctx, tt.order)

			if tt.errType != nil {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tt.errType)
			} else {
				require.NoError(s.T(), err)

				updatedOrder, getErr := s.repo.GetOrder(s.ctx, tt.order.OrderUUID.String())
				require.NoError(s.T(), getErr)
				require.NotNil(s.T(), updatedOrder)
				require.Equal(s.T(), tt.order.OrderUUID, updatedOrder.OrderUUID)
				require.Equal(s.T(), tt.order.UserUUID, updatedOrder.UserUUID)
				require.Equal(s.T(), tt.order.TotalPrice, updatedOrder.TotalPrice)
				require.Equal(s.T(), tt.order.Status, updatedOrder.Status)
				require.Equal(s.T(), len(tt.order.PartUUIDs), len(updatedOrder.PartUUIDs))
				if tt.order.TransactionUUID != nil {
					require.NotNil(s.T(), updatedOrder.TransactionUUID)
					require.Equal(s.T(), *tt.order.TransactionUUID, *updatedOrder.TransactionUUID)
				}
				if tt.order.PaymentMethod != nil {
					require.NotNil(s.T(), updatedOrder.PaymentMethod)
					require.Equal(s.T(), *tt.order.PaymentMethod, *updatedOrder.PaymentMethod)
				}
			}
		})
	}
}
