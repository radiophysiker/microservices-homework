package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/radiophysiker/microservices-homework/iam/internal/model"
	"github.com/radiophysiker/microservices-homework/iam/internal/repository/converter"
)

const (
	uniqueViolationCode = "23505"
)

// Create сохраняет пользователя в PostgreSQL.
func (r *Repository) Create(ctx context.Context, user *model.User) error {
	repoUser, err := converter.ToRepoUser(user)
	if err != nil {
		return fmt.Errorf("convert user to repository model: %w", err)
	}

	const query = `
INSERT INTO users (uuid, login, email, password_hash, notification_methods, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);
`

	_, err = r.pool.Exec(
		ctx,
		query,
		repoUser.UUID,
		repoUser.Login,
		repoUser.Email,
		repoUser.PasswordHash,
		repoUser.NotificationMethods,
		repoUser.CreatedAt,
		repoUser.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return model.NewErrUserAlreadyExists(user.UUID)
		}

		return fmt.Errorf("exec insert user: %w", err)
	}

	return nil
}
