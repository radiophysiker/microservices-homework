package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/radiophysiker/microservices-homework/payment/internal/model"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
)

// PayOrder проводит оплату заказа
func (a *API) PayOrder(ctx context.Context, req *pb.PayOrderRequest) (*pb.PayOrderResponse, error) {
	transactionUUID, err := a.paymentService.PayOrder(ctx, req.UserUuid, req.OrderUuid, req.PaymentMethod)
	if err != nil {
		if errors.Is(err, model.ErrInvalidPaymentRequest) {
			return nil, status.Error(codes.InvalidArgument, "invalid payment request")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
