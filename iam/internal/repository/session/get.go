package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/radiophysiker/microservices-homework/iam/internal/model"
	"github.com/radiophysiker/microservices-homework/iam/internal/repository/converter"
	repoModel "github.com/radiophysiker/microservices-homework/iam/internal/repository/model"
)

// Get возвращает сессию по ее UUID.
func (r *Repository) Get(ctx context.Context, sessionUUID string) (*model.Session, error) {
	key := sessionKey(sessionUUID)

	data, err := r.client.Get(ctx, key)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return nil, model.ErrInvalidCredentials
		}

		return nil, fmt.Errorf("get session from redis: %w", err)
	}

	var repoSession repoModel.Session
	if err := json.Unmarshal(data, &repoSession); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}

	session := converter.ToServiceSession(&repoSession)
	if session == nil {
		return nil, model.ErrInvalidCredentials
	}

	return session, nil
}
