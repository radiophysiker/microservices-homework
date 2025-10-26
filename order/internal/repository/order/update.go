package order

import (
	"context"
	"fmt"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	"github.com/radiophysiker/microservices-homework/order/internal/repository/converter"
)

// UpdateOrder обновляет заказ
func (r *Repository) UpdateOrder(ctx context.Context, order *model.Order) error {
	repoOrder := converter.ToRepoOrder(order)

	query := `
		UPDATE orders 
		SET user_uuid = $2, part_uuids = $3, total_price = $4, transaction_uuid = $5, payment_method = $6, status = $7, updated_at = NOW()
		WHERE uuid = $1
	`

	var paymentMethodStr *string

	if repoOrder.PaymentMethod != nil {
		str := repoOrder.PaymentMethod.String()
		paymentMethodStr = &str
	}

	result, err := r.pool.Exec(ctx, query,
		repoOrder.OrderUUID,
		repoOrder.UserUUID,
		repoOrder.PartUUIDs,
		repoOrder.TotalPrice,
		repoOrder.TransactionUUID,
		paymentMethodStr,
		repoOrder.Status.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}
