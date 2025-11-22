package session

import (
	"context"
	"fmt"
)

// AddSessionToUserSet добавляет идентификатор сессии в множество пользователя.
func (r *Repository) AddSessionToUserSet(ctx context.Context, userUUID, sessionUUID string) error {
	key := userSessionsKey(userUUID)

	if err := r.client.SAdd(ctx, key, sessionUUID); err != nil {
		return fmt.Errorf("add session to user set: %w", err)
	}

	if err := r.client.Expire(ctx, key, r.ttl); err != nil {
		return fmt.Errorf("set TTL for user session set: %w", err)
	}

	return nil
}
