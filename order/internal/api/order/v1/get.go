package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/radiophysiker/microservices-homework/order/internal/converter"
	"github.com/radiophysiker/microservices-homework/order/internal/model"
	orderpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/order/v1"
)

// GetOrder возвращает заказ по UUID
func (a *API) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	orderUUID, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order UUID: %v", err)
	}

	order, err := a.orderService.GetOrder(ctx, orderUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidOrderData):
			return nil, status.Errorf(codes.InvalidArgument, "invalid order data: %v", err)
		case errors.Is(err, model.ErrOrderNotFound):
			return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
		default:
			return nil, status.Errorf(codes.Internal, "failed to get order: %v", err)
		}
	}

	return converter.ToProtoOrder(order), nil
}
