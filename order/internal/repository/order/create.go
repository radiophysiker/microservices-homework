package order

import (
	"context"
	"fmt"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	"github.com/radiophysiker/microservices-homework/order/internal/repository/converter"
)

// CreateOrder создает новый заказ
func (r *Repository) CreateOrder(ctx context.Context, order *model.Order) error {
	repoOrder := converter.ToRepoOrder(order)

	query := `
		INSERT INTO orders (uuid, user_uuid, part_uuids, total_price, transaction_uuid, payment_method, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	var paymentMethodStr *string

	if repoOrder.PaymentMethod != nil {
		str := repoOrder.PaymentMethod.String()
		paymentMethodStr = &str
	}

	_, err := r.pool.Exec(ctx, query,
		repoOrder.OrderUUID,
		repoOrder.UserUUID,
		repoOrder.PartUUIDs,
		repoOrder.TotalPrice,
		repoOrder.TransactionUUID,
		paymentMethodStr,
		repoOrder.Status.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}
