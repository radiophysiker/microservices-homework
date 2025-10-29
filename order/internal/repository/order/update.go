package order

import (
	"context"
	"errors"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/radiophysiker/microservices-homework/order/internal/model"
	"github.com/radiophysiker/microservices-homework/order/internal/repository/converter"
	repoModel "github.com/radiophysiker/microservices-homework/order/internal/repository/model"
)

// UpdateOrder обновляет заказ и возвращает актуальное состояние
func (r *Repository) UpdateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	repoOrder := converter.ToRepoOrder(order)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer r.rollbackTx(ctx, tx)

	if err := r.updateOrderTable(ctx, tx, repoOrder); err != nil {
		return nil, err
	}

	if err := r.upsertOrderItems(ctx, tx, repoOrder); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	updated, getErr := r.GetOrder(ctx, repoOrder.OrderUUID.String())
	if getErr != nil {
		return nil, fmt.Errorf("failed to fetch updated order: %w", getErr)
	}

	return updated, nil
}

// updateOrderTable обновляет данные в таблице orders
func (r *Repository) updateOrderTable(ctx context.Context, tx pgx.Tx, repoOrder *repoModel.Order) error {
	var paymentMethodStr *string

	if repoOrder.PaymentMethod != nil {
		str := repoOrder.PaymentMethod.String()
		paymentMethodStr = &str
	}

	query, args, err := sq.Update("orders").
		Set("user_uuid", repoOrder.UserUUID).
		Set("total_price", repoOrder.TotalPrice).
		Set("transaction_uuid", repoOrder.TransactionUUID).
		Set("payment_method", paymentMethodStr).
		Set("status", repoOrder.Status.String()).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"uuid": repoOrder.OrderUUID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update order query: %w", err)
	}

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}

// upsertOrderItems вставляет новые или обновляет существующие items,
// а также удаляет items, которых нет в новом списке
func (r *Repository) upsertOrderItems(ctx context.Context, tx pgx.Tx, repoOrder *repoModel.Order) error {
	if len(repoOrder.Items) == 0 {
		return r.deleteAllOrderItems(ctx, tx, repoOrder.OrderUUID)
	}

	// 1. UPSERT существующих и новых items
	builder := sq.Insert("order_items").
		Columns("order_uuid", "part_uuid", "quantity").
		PlaceholderFormat(sq.Dollar).
		Suffix("ON CONFLICT (order_uuid, part_uuid) DO UPDATE SET quantity = EXCLUDED.quantity")

	partUUIDs := make([]uuid.UUID, 0, len(repoOrder.Items))

	for _, item := range repoOrder.Items {
		builder = builder.Values(
			repoOrder.OrderUUID,
			item.PartUUID,
			item.Quantity,
		)

		partUUIDs = append(partUUIDs, item.PartUUID)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build upsert order items query: %w", err)
	}

	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to upsert order items: %w", err)
	}

	if err := r.deleteObsoleteOrderItems(ctx, tx, repoOrder.OrderUUID, partUUIDs); err != nil {
		return err
	}

	return nil
}

// deleteObsoleteOrderItems удаляет items, которых нет в списке актуальных partUUIDs
func (r *Repository) deleteObsoleteOrderItems(ctx context.Context, tx pgx.Tx, orderUUID uuid.UUID, keepPartUUIDs []uuid.UUID) error {
	query, args, err := sq.Delete("order_items").
		Where(sq.Eq{"order_uuid": orderUUID}).
		Where(sq.NotEq{"part_uuid": keepPartUUIDs}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete obsolete items query: %w", err)
	}

	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to delete obsolete items: %w", err)
	}

	return nil
}

// deleteAllOrderItems удаляет все items заказа
func (r *Repository) deleteAllOrderItems(ctx context.Context, tx pgx.Tx, orderUUID uuid.UUID) error {
	query, args, err := sq.Delete("order_items").
		Where(sq.Eq{"order_uuid": orderUUID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete all order items query: %w", err)
	}

	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to delete all order items: %w", err)
	}

	return nil
}

// rollbackTx откатывает транзакцию, если она не была закоммичена
func (r *Repository) rollbackTx(ctx context.Context, tx pgx.Tx) {
	if tx == nil {
		return
	}

	if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		log.Printf("failed to rollback transaction: %v\n", err)
	}
}
