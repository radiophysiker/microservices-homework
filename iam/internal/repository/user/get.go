package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/radiophysiker/microservices-homework/iam/internal/model"
	"github.com/radiophysiker/microservices-homework/iam/internal/repository/converter"
	repoModel "github.com/radiophysiker/microservices-homework/iam/internal/repository/model"
)

// GetByUUID возвращает пользователя по UUID.
func (r *Repository) GetByUUID(ctx context.Context, userUUID string) (*model.User, error) {
	const query = `
SELECT uuid, login, email, password_hash, notification_methods, created_at, updated_at
FROM users
WHERE uuid = $1
LIMIT 1;
`

	return r.get(ctx, query, userUUID)
}

// GetByLogin возвращает пользователя по логину.
func (r *Repository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	const query = `
SELECT uuid, login, email, password_hash, notification_methods, created_at, updated_at
FROM users
WHERE login = $1
LIMIT 1;
`

	return r.get(ctx, query, login)
}

// GetByEmail возвращает пользователя по email.
func (r *Repository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	const query = `
SELECT uuid, login, email, password_hash, notification_methods, created_at, updated_at
FROM users
WHERE LOWER(email) = LOWER($1)
LIMIT 1;
`

	return r.get(ctx, query, email)
}

// get выполняет запрос к базе данных и возвращает пользователя.
// Используется внутренними методами GetByUUID, GetByLogin и GetByEmail.
func (r *Repository) get(ctx context.Context, query string, arg any) (*model.User, error) {
	var repoUser repoModel.User

	err := r.pool.
		QueryRow(ctx, query, arg).
		Scan(
			&repoUser.UUID,
			&repoUser.Login,
			&repoUser.Email,
			&repoUser.PasswordHash,
			&repoUser.NotificationMethods,
			&repoUser.CreatedAt,
			&repoUser.UpdatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}

		return nil, fmt.Errorf("query user: %w", err)
	}

	user, err := converter.ToServiceUser(&repoUser)
	if err != nil {
		return nil, fmt.Errorf("convert user: %w", err)
	}

	return user, nil
}
