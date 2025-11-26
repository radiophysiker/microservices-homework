package v1

import (
	"context"
	"fmt"

	grpcMiddleware "github.com/radiophysiker/microservices-homework/platform/pkg/middleware/grpc"
	paymentpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/payment/v1"
)

// Client реализует интерфейс PaymentClient
type Client struct {
	paymentClient paymentpb.PaymentServiceClient
}

// NewClient создает новый экземпляр Client
func NewClient(paymentClient paymentpb.PaymentServiceClient) *Client {
	return &Client{
		paymentClient: paymentClient,
	}
}

// PayOrder проводит оплату заказа
func (c *Client) PayOrder(ctx context.Context, userUUID, orderUUID string, paymentMethod paymentpb.PaymentMethod) (string, error) {
	ctx = grpcMiddleware.ForwardSessionUUIDToGRPC(ctx)

	resp, err := c.paymentClient.PayOrder(ctx, &paymentpb.PayOrderRequest{
		UserUuid:      userUUID,
		OrderUuid:     orderUUID,
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return "", fmt.Errorf("failed to pay order: %w", err)
	}

	return resp.GetTransactionUuid(), nil
}
