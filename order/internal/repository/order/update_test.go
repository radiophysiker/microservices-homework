package order

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
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
			order: s.createTestOrder(s.testOrderUUID, 250.00, model.StatusPaid),
			setupFunc: func() {
				existingOrder := s.createTestOrder(s.testOrderUUID, 100.00, model.StatusPendingPayment)
				s.repo.On("CreateOrder", s.ctx, existingOrder).Return(nil).Once()
				err := s.repo.CreateOrder(s.ctx, existingOrder)
				require.NoError(s.T(), err)
			},
			errType: nil,
		},
		{
			name:  "success_change_status",
			order: s.createTestOrder(s.testOrderUUID, 100.00, model.StatusCancelled),
			setupFunc: func() {
				existingOrder := s.createTestOrder(s.testOrderUUID, 100.00, model.StatusPendingPayment)
				s.repo.On("CreateOrder", s.ctx, existingOrder).Return(nil).Once()
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
				model.StatusPendingPayment,
			),
			setupFunc: func() {
				existingOrder := s.createTestOrder(s.testOrderUUID, 100.00, model.StatusPendingPayment)
				s.repo.On("CreateOrder", s.ctx, existingOrder).Return(nil).Once()
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
				model.StatusPaid,
				ptrUUID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440008")),
				ptrPaymentMethod(model.PaymentMethodCard),
			),
			setupFunc: func() {
				existingOrder := s.createTestOrder(s.testOrderUUID, 100.00, model.StatusPendingPayment)
				s.repo.On("CreateOrder", s.ctx, existingOrder).Return(nil).Once()
				err := s.repo.CreateOrder(s.ctx, existingOrder)
				require.NoError(s.T(), err)
			},
			errType: nil,
		},
		{
			name:      "not_found",
			order:     s.createTestOrder(uuid.MustParse("550e8400-e29b-41d4-a716-999999999999"), 100.00, model.StatusPendingPayment),
			setupFunc: nil,
			errType:   model.ErrOrderNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.setupFunc != nil {
				tt.setupFunc()
			}

			if tt.errType != nil {
				s.repo.On("UpdateOrder", s.ctx, tt.order).Return(tt.errType).Once()
			} else {
				s.repo.On("UpdateOrder", s.ctx, tt.order).Return(nil).Once()
				s.repo.On("GetOrder", s.ctx, tt.order.OrderUUID.String()).Return(tt.order, nil).Once()
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
