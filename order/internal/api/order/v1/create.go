package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	orderpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/order/v1"
)

// CreateOrder создает новый заказ (gRPC)
func (a *API) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	userUUID, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user UUID: %v", err)
	}

	partUUIDs := make([]uuid.UUID, len(req.GetPartUuids()))

	for i, partUUIDStr := range req.GetPartUuids() {
		partUUID, err := uuid.Parse(partUUIDStr)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid part UUID at index %d: %v", i, err)
		}

		partUUIDs[i] = partUUID
	}

	order, err := a.orderService.CreateOrder(ctx, userUUID, partUUIDs)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidOrderData):
			return nil, status.Errorf(codes.InvalidArgument, "invalid order data: %v", err)
		case errors.Is(err, model.ErrInventoryServiceUnavailable):
			return nil, status.Errorf(codes.Unavailable, "inventory service unavailable: %v", err)
		default:
			return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
		}
	}

	return &orderpb.CreateOrderResponse{
		OrderUuid:  order.OrderUUID.String(),
		TotalPrice: order.TotalPrice,
	}, nil
}
