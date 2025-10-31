package order

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	"github.com/radiophysiker/microservices-homework/order/internal/repository/converter"
	repoModel "github.com/radiophysiker/microservices-homework/order/internal/repository/model"
)

// GetOrder возвращает заказ по UUID
func (r *Repository) GetOrder(ctx context.Context, orderUUID string) (*model.Order, error) {
	_, err := uuid.Parse(orderUUID)
	if err != nil {
		return nil, model.NewInvalidOrderDataError(orderUUID)
	}

	var (
		repoOrder        repoModel.Order
		paymentMethodStr *string
		statusStr        string
	)

	builder := sq.
		Select(
			"uuid",
			"user_uuid",
			"total_price",
			"transaction_uuid",
			"payment_method",
			"status",
		).
		From("orders").
		Where(sq.Eq{"uuid": orderUUID}).
		PlaceholderFormat(sq.Dollar)

	sql, args, buildErr := builder.ToSql()
	if buildErr != nil {
		return nil, fmt.Errorf("failed to build get order query: %w", buildErr)
	}

	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&repoOrder.OrderUUID,
		&repoOrder.UserUUID,
		&repoOrder.TotalPrice,
		&repoOrder.TransactionUUID,
		&paymentMethodStr,
		&statusStr,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrOrderNotFound
		}

		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if paymentMethodStr != nil {
		repoOrder.PaymentMethod = converter.StringToPaymentMethod(*paymentMethodStr)
	}

	repoOrder.Status = converter.StringToOrderStatus(statusStr)

	itemsSQL, itemsArgs, buildItemsErr := sq.Select("part_uuid", "quantity").
		From("order_items").
		Where(sq.Eq{"order_uuid": repoOrder.OrderUUID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if buildItemsErr != nil {
		return nil, fmt.Errorf("failed to build get order items query: %w", buildItemsErr)
	}

	rows, qErr := r.pool.Query(ctx, itemsSQL, itemsArgs...)
	if qErr != nil {
		return nil, fmt.Errorf("failed to get order items: %w", qErr)
	}
	defer rows.Close()

	for rows.Next() {
		var it repoModel.OrderItem
		if err = rows.Scan(&it.PartUUID, &it.Quantity); err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		repoOrder.Items = append(repoOrder.Items, it)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to iterate order items: %w", rows.Err())
	}

	return converter.ToServiceOrder(&repoOrder), nil
}
