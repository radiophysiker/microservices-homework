package order

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	"github.com/radiophysiker/microservices-homework/order/internal/repository/converter"
)

// CreateOrder создает новый заказ
func (r *Repository) CreateOrder(ctx context.Context, order *model.Order) error {
	repoOrder := converter.ToRepoOrder(order)

	var paymentMethodStr *string
	if repoOrder.PaymentMethod != nil {
		str := repoOrder.PaymentMethod.String()
		paymentMethodStr = &str
	}

	builder := sq.Insert("orders").
		Columns("uuid", "user_uuid", "total_price", "transaction_uuid", "payment_method", "status").
		Values(
			repoOrder.OrderUUID,
			repoOrder.UserUUID,
			repoOrder.TotalPrice,
			repoOrder.TransactionUUID,
			paymentMethodStr,
			repoOrder.Status.String(),
		).PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build create order query: %w", err)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	defer func() {
		if tx != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				log.Printf("failed to rollback tx: %v", err)
			}
		}
	}()

	if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	if len(repoOrder.Items) > 0 {
		itemsInsert := sq.Insert("order_items").
			Columns("order_uuid", "part_uuid", "quantity").
			PlaceholderFormat(sq.Dollar)

		for _, it := range repoOrder.Items {
			itemsInsert = itemsInsert.Values(repoOrder.OrderUUID, it.PartUUID, it.Quantity)
		}

		itemSQL, itemArgs, buildErr := itemsInsert.ToSql()
		if buildErr != nil {
			return fmt.Errorf("failed to build insert order items query: %w", buildErr)
		}

		if _, err = tx.Exec(ctx, itemSQL, itemArgs...); err != nil {
			return fmt.Errorf("failed to insert order items: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	tx = nil

	return nil
}
