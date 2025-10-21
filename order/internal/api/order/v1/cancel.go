package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	orderpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/order/v1"
)

// CancelOrder отменяет заказ
func (a *API) CancelOrder(ctx context.Context, req *orderpb.CancelOrderRequest) (*emptypb.Empty, error) {
	orderUUID, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order UUID: %v", err)
	}

	_, err = a.orderService.CancelOrder(ctx, orderUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidOrderData):
			return nil, status.Errorf(codes.InvalidArgument, "invalid order data: %v", err)
		case errors.Is(err, model.ErrOrderNotFound):
			return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
		case errors.Is(err, model.ErrOrderCannotBeCancelled):
			return nil, status.Errorf(codes.FailedPrecondition, "order cannot be cancelled: %v", err)
		default:
			return nil, status.Errorf(codes.Internal, "failed to cancel order: %v", err)
		}
	}

	return &emptypb.Empty{}, nil
}
