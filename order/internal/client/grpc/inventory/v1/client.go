package v1

import (
	"context"
	"fmt"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	inventorypb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
)

// Client реализует интерфейс InventoryClient
type Client struct {
	inventoryClient inventorypb.InventoryServiceClient
}

// NewClient создает новый экземпляр Client
func NewClient(inventoryClient inventorypb.InventoryServiceClient) *Client {
	return &Client{
		inventoryClient: inventoryClient,
	}
}

// ListParts возвращает список деталей по UUID
func (c *Client) ListParts(ctx context.Context, partUUIDs []string) ([]*model.Part, error) {
	resp, err := c.inventoryClient.ListParts(ctx, &inventorypb.ListPartsRequest{
		Filter: &inventorypb.PartsFilter{
			Uuids: partUUIDs,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list parts: %w", err)
	}

	return model.ToServiceParts(resp.Parts), nil
}
