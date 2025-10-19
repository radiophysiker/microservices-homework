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
)

func (s *ServiceTestSuite) TestCreateOrder() {
	tests := []struct {
		name      string
		userUUID  uuid.UUID
		partUUIDs []uuid.UUID
		setupMock func(*repomocks.MockOrderRepository, *clientmocks.MockInventoryClient, *clientmocks.MockPaymentClient)
		wantOrder *model.Order
		checkErr  func(err error)
	}{
		{
			name:      "success",
			userUUID:  uuid.New(),
			partUUIDs: []uuid.UUID{uuid.New(), uuid.New()},
			setupMock: func(repo *repomocks.MockOrderRepository, inv *clientmocks.MockInventoryClient, pay *clientmocks.MockPaymentClient) {
				parts := []*model.Part{{UUID: uuid.New().String(), Price: 10}, {UUID: uuid.New().String(), Price: 25}}
				inv.EXPECT().ListParts(s.ctx, mock.AnythingOfType("[]string")).Return(parts, nil).Once()
				repo.EXPECT().CreateOrder(s.ctx, mock.AnythingOfType("*model.Order")).Return(nil).Once()
			},
			wantOrder: &model.Order{
				Status: orderv1.OrderStatusPENDINGPAYMENT,
			},
		},
		{
			name:      "inventory_error",
			userUUID:  uuid.New(),
			partUUIDs: []uuid.UUID{uuid.New()},
			setupMock: func(repo *repomocks.MockOrderRepository, inv *clientmocks.MockInventoryClient, pay *clientmocks.MockPaymentClient) {
				inv.EXPECT().ListParts(s.ctx, mock.AnythingOfType("[]string")).Return(nil, errors.New("inventory service down")).Once()
			},
			wantOrder: nil,
			checkErr: func(err error) {
				assert.Error(s.T(), err)
				assert.Contains(s.T(), err.Error(), "inventory service unavailable")
				assert.Contains(s.T(), err.Error(), "inventory service down")
			},
		},
		{
			name:      "repository_error",
			userUUID:  uuid.New(),
			partUUIDs: []uuid.UUID{uuid.New()},
			setupMock: func(repo *repomocks.MockOrderRepository, inv *clientmocks.MockInventoryClient, pay *clientmocks.MockPaymentClient) {
				parts := []*model.Part{{UUID: uuid.New().String(), Price: 10}}
				inv.EXPECT().ListParts(s.ctx, mock.AnythingOfType("[]string")).Return(parts, nil).Once()
				repo.EXPECT().CreateOrder(s.ctx, mock.AnythingOfType("*model.Order")).Return(errors.New("database error")).Once()
			},
			wantOrder: nil,
			checkErr: func(err error) {
				assert.Error(s.T(), err)
				assert.Contains(s.T(), err.Error(), "failed to create order")
				assert.Contains(s.T(), err.Error(), "database error")
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock(s.repo, s.inventoryClient, s.paymentClient)

			got, err := s.service.CreateOrder(s.ctx, tt.userUUID, tt.partUUIDs)

			if tt.checkErr != nil {
				tt.checkErr(err)

				if err != nil {
					require.Nil(s.T(), got)
				}
			} else {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), got)
				require.Equal(s.T(), tt.wantOrder.Status, got.Status)
				require.Equal(s.T(), tt.userUUID, got.UserUUID)
			}
		})
	}
}
