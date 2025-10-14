package v1

import (
	"context"
	"errors"

	"github.com/radiophysiker/microservices-homework/payment/internal/model"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PayOrder проводит оплату заказа
func (a *API) PayOrder(ctx context.Context, req *pb.PayOrderRequest) (*pb.PayOrderResponse, error) {
	transactionUUID, err := a.paymentService.PayOrder(ctx, req.UserUuid, req.OrderUuid, req.PaymentMethod)
	if err != nil {
		if errors.Is(err, model.ErrInvalidPaymentRequest) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid payment request: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "payment processing failed: %v", err)
	}

	return &pb.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
