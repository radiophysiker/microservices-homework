package order

import (
	"github.com/stretchr/testify/require"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
)

func (s *RepositoryTestSuite) TestGetOrder() {
	existingOrder := &model.Order{
		OrderUUID:  s.testOrderUUID,
		UserUUID:   s.testUserUUID,
		Items:      []model.OrderItem{{PartUUID: s.testPartUUID, Quantity: 1}},
		TotalPrice: 100.50,
		Status:     model.StatusPendingPayment,
	}

	s.repo.On("CreateOrder", s.ctx, existingOrder).Return(nil).Once()

	err := s.repo.CreateOrder(s.ctx, existingOrder)
	require.NoError(s.T(), err, "failed to setup test data")

	tests := []struct {
		name    string
		uuid    string
		wantErr bool
		errType error
	}{
		{
			name:    "success",
			uuid:    s.testOrderUUID.String(),
			wantErr: false,
			errType: nil,
		},
		{
			name:    "not_found",
			uuid:    "550e8400-e29b-41d4-a716-999999999999",
			wantErr: true,
			errType: model.ErrOrderNotFound,
		},
		{
			name:    "invalid_uuid_format",
			uuid:    "not-a-valid-uuid",
			wantErr: true,
			errType: model.ErrInvalidOrderData,
		},
		{
			name:    "empty_uuid",
			uuid:    "",
			wantErr: true,
			errType: model.ErrInvalidOrderData,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.wantErr {
				s.repo.On("GetOrder", s.ctx, tt.uuid).Return(nil, tt.errType).Once()
			} else {
				s.repo.On("GetOrder", s.ctx, tt.uuid).Return(existingOrder, nil).Once()
			}

			order, err := s.repo.GetOrder(s.ctx, tt.uuid)

			if tt.wantErr {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tt.errType)
				require.Nil(s.T(), order)
			} else {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), order)
				require.Equal(s.T(), tt.uuid, order.OrderUUID.String())
			}
		})
	}
}
