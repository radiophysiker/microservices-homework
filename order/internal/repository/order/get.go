package order

import (
	"context"
	"errors"
	"fmt"

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

	query := `
		SELECT uuid, user_uuid, part_uuids, total_price, transaction_uuid, payment_method, status
		FROM orders
		WHERE uuid = $1
	`

	var repoOrder repoModel.Order

	var paymentMethodStr *string

	var statusStr string

	err = r.pool.QueryRow(ctx, query, orderUUID).Scan(
		&repoOrder.OrderUUID,
		&repoOrder.UserUUID,
		&repoOrder.PartUUIDs,
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

	repoOrder.Status = converter.StringToStatus(statusStr)

	return converter.ToServiceOrder(&repoOrder), nil
}
