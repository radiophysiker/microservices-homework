package order

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	clientmocks "github.com/radiophysiker/microservices-homework/order/internal/client/grpc/mocks"
	"github.com/radiophysiker/microservices-homework/order/internal/model"
	repomocks "github.com/radiophysiker/microservices-homework/order/internal/repository/mocks"
	orderv1 "github.com/radiophysiker/microservices-homework/shared/pkg/openapi/order/v1"
	paymentpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
)

func (s *ServiceTestSuite) TestPayOrder() {
	tests := []struct {
		name          string
		orderUUID     uuid.UUID
		paymentMethod orderv1.PaymentMethod
		setupMock     func(*repomocks.MockOrderRepository, *clientmocks.MockPaymentClient)
		wantOrder     *model.Order
		checkErr      func(err error)
	}{
		{
			name:          "success",
			orderUUID:     uuid.New(),
			paymentMethod: orderv1.PaymentMethodCARD,
			setupMock: func(repo *repomocks.MockOrderRepository, pay *clientmocks.MockPaymentClient) {
				order := &model.Order{
					OrderUUID:  uuid.New(),
					UserUUID:   uuid.New(),
					Status:     orderv1.OrderStatusPENDINGPAYMENT,
					TotalPrice: 100,
				}
				repo.EXPECT().GetOrder(s.ctx, mock.AnythingOfType("string")).Return(order, nil).Once()
				pay.EXPECT().PayOrder(s.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.MatchedBy(func(pm paymentpb.PaymentMethod) bool { return true })).Return("550e8400-e29b-41d4-a716-446655440000", nil).Once()
				repo.EXPECT().UpdateOrder(s.ctx, mock.AnythingOfType("*model.Order")).Return(nil).Once()
			},
			wantOrder: &model.Order{
				Status: orderv1.OrderStatusPAID,
			},
		},
		{
			name:          "get_order_error",
			orderUUID:     uuid.New(),
			paymentMethod: orderv1.PaymentMethodCARD,
			setupMock: func(repo *repomocks.MockOrderRepository, pay *clientmocks.MockPaymentClient) {
				repo.EXPECT().GetOrder(s.ctx, mock.AnythingOfType("string")).Return((*model.Order)(nil), errors.New("order not found")).Once()
			},
			wantOrder: nil,
			checkErr: func(err error) {
				assert.Error(s.T(), err)
				assert.Contains(s.T(), err.Error(), "failed to get order")
				assert.Contains(s.T(), err.Error(), "order not found")
			},
		},
		{
			name:          "payment_error",
			orderUUID:     uuid.New(),
			paymentMethod: orderv1.PaymentMethodCARD,
			setupMock: func(repo *repomocks.MockOrderRepository, pay *clientmocks.MockPaymentClient) {
				order := &model.Order{
					OrderUUID:  uuid.New(),
					UserUUID:   uuid.New(),
					Status:     orderv1.OrderStatusPENDINGPAYMENT,
					TotalPrice: 100,
				}
				repo.EXPECT().GetOrder(s.ctx, mock.AnythingOfType("string")).Return(order, nil).Once()
				pay.EXPECT().PayOrder(s.ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.MatchedBy(func(pm paymentpb.PaymentMethod) bool { return true })).Return("", errors.New("payment failed")).Once()
			},
			wantOrder: nil,
			checkErr: func(err error) {
				assert.Error(s.T(), err)
				assert.Contains(s.T(), err.Error(), "payment service unavailable")
				assert.Contains(s.T(), err.Error(), "payment failed")
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock(s.repo, s.paymentClient)

			got, err := s.service.PayOrder(s.ctx, tt.orderUUID, tt.paymentMethod)

			if tt.checkErr != nil {
				tt.checkErr(err)
			} else {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), got)
				require.Equal(s.T(), tt.wantOrder.Status, got.Status)
			}
		})
	}
}
