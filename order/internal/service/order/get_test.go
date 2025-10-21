package order

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	repomocks "github.com/radiophysiker/microservices-homework/order/internal/repository/mocks"
)

func (s *ServiceTestSuite) TestGetOrder() {
	tests := []struct {
		name      string
		orderUUID uuid.UUID
		setupMock func(*repomocks.MockOrderRepository)
		wantOrder *model.Order
		checkErr  func(err error)
	}{
		{
			name:      "success",
			orderUUID: uuid.New(),
			setupMock: func(repo *repomocks.MockOrderRepository) {
				want := &model.Order{
					OrderUUID: uuid.New(),
					UserUUID:  uuid.New(),
					Status:    model.StatusPendingPayment,
				}
				repo.EXPECT().GetOrder(s.ctx, mock.AnythingOfType("string")).Return(want, nil).Once()
			},
			wantOrder: &model.Order{
				Status: model.StatusPendingPayment,
			},
		},
		{
			name:      "repository_error",
			orderUUID: uuid.New(),
			setupMock: func(repo *repomocks.MockOrderRepository) {
				repo.EXPECT().GetOrder(s.ctx, mock.AnythingOfType("string")).Return((*model.Order)(nil), errors.New("database error")).Once()
			},
			wantOrder: nil,
			checkErr: func(err error) {
				assert.Error(s.T(), err)
				assert.Contains(s.T(), err.Error(), "failed to get order")
				assert.Contains(s.T(), err.Error(), "database error")
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock(s.repo)

			got, err := s.service.GetOrder(s.ctx, tt.orderUUID)

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
