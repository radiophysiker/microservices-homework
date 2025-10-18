package payment

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/radiophysiker/microservices-homework/payment/internal/model"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	svc *Service
	ctx context.Context
}

func (s *ServiceSuite) SetupTest() {
	s.svc = NewService()
	s.ctx = context.Background()
}

func (s *ServiceSuite) TestPayOrder() {
	tests := []struct {
		name          string
		userUUID      string
		orderUUID     string
		method        pb.PaymentMethod
		wantErr       error
		wantErrSubstr string
	}{
		{name: "success_card", userUUID: "user-123", orderUUID: "order-456", method: pb.PaymentMethod_PAYMENT_METHOD_CARD},
		{name: "success_sbp", userUUID: "user-789", orderUUID: "order-012", method: pb.PaymentMethod_PAYMENT_METHOD_SBP},
		{name: "success_credit_card", userUUID: "user-111", orderUUID: "order-222", method: pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD},
		{name: "success_investor_money", userUUID: "user-333", orderUUID: "order-444", method: pb.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY},

		{name: "invalid_user_uuid_empty", userUUID: "", orderUUID: "order-456", method: pb.PaymentMethod_PAYMENT_METHOD_CARD, wantErr: model.ErrInvalidPaymentRequest},
		{name: "invalid_order_uuid_empty", userUUID: "user-123", orderUUID: "", method: pb.PaymentMethod_PAYMENT_METHOD_CARD, wantErr: model.ErrInvalidPaymentRequest},
		{name: "invalid_payment_method_unspecified", userUUID: "user-123", orderUUID: "order-456", method: pb.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, wantErr: model.ErrInvalidPaymentRequest, wantErrSubstr: "unspecified payment method"},
	}

	for _, tt := range tests {
		tt := tt
		s.Run(tt.name, func() {
			id, err := s.svc.PayOrder(s.ctx, tt.userUUID, tt.orderUUID, tt.method)

			if tt.wantErr != nil {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tt.wantErr)
				if tt.wantErrSubstr != "" {
					require.Contains(s.T(), err.Error(), tt.wantErrSubstr)
				}
				return
			}

			require.NoError(s.T(), err)
			require.NotEmpty(s.T(), id)
			_, perr := uuid.Parse(id)
			require.NoError(s.T(), perr)
		})
	}
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
