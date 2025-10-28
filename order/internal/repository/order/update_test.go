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
				s.repo.On("UpdateOrder", s.ctx, tt.order).Return((*model.Order)(nil), tt.errType).Once()
			} else {
				s.repo.On("UpdateOrder", s.ctx, tt.order).Return(tt.order, nil).Once()
			}

			updated, err := s.repo.UpdateOrder(s.ctx, tt.order)

			if tt.errType != nil {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tt.errType)
			} else {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), updated)
				require.Equal(s.T(), tt.order.OrderUUID, updated.OrderUUID)
				require.Equal(s.T(), tt.order.UserUUID, updated.UserUUID)
				require.Equal(s.T(), tt.order.TotalPrice, updated.TotalPrice)
				require.Equal(s.T(), tt.order.Status, updated.Status)
				require.Equal(s.T(), len(tt.order.Items), len(updated.Items))

				if tt.order.TransactionUUID != nil {
					require.NotNil(s.T(), updated.TransactionUUID)
					require.Equal(s.T(), *tt.order.TransactionUUID, *updated.TransactionUUID)
				}

				if tt.order.PaymentMethod != nil {
					require.NotNil(s.T(), updated.PaymentMethod)
					require.Equal(s.T(), *tt.order.PaymentMethod, *updated.PaymentMethod)
				}
			}
		})
	}
}
