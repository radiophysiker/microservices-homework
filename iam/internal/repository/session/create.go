package session

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/radiophysiker/microservices-homework/iam/internal/model"
	"github.com/radiophysiker/microservices-homework/iam/internal/repository/converter"
)

// Create сохраняет сессию в Redis с TTL.
func (r *Repository) Create(ctx context.Context, session *model.Session) error {
	repoSession := converter.ToRepoSession(session)
	if repoSession == nil {
		return fmt.Errorf("session is nil")
	}

	payload, err := json.Marshal(repoSession)
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}

	key := sessionKey(repoSession.UUID)
	if err := r.client.SetWithTTL(ctx, key, payload, r.ttl); err != nil {
		return fmt.Errorf("set session in redis: %w", err)
	}

	return nil
}
