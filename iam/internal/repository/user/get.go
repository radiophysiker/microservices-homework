package user

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/radiophysiker/microservices-homework/iam/internal/model"
	"github.com/radiophysiker/microservices-homework/iam/internal/repository/converter"
	repoModel "github.com/radiophysiker/microservices-homework/iam/internal/repository/model"
)

// GetByUUID возвращает пользователя по UUID.
func (r *Repository) GetByUUID(ctx context.Context, userUUID string) (*model.User, error) {
	query, args, err := buildGetUserQuery(sq.Eq{"uuid": userUUID})
	if err != nil {
		return nil, fmt.Errorf("build get user by uuid query: %w", err)
	}

	return r.get(ctx, query, args...)
}

// GetByLogin возвращает пользователя по логину.
func (r *Repository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	query, args, err := buildGetUserQuery(sq.Eq{"login": login})
	if err != nil {
		return nil, fmt.Errorf("build get user by login query: %w", err)
	}

	return r.get(ctx, query, args...)
}

// GetByEmail возвращает пользователя по email.
func (r *Repository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query, args, err := buildGetUserQuery(sq.Eq{"email": email})
	if err != nil {
		return nil, fmt.Errorf("build get user by email query: %w", err)
	}

	return r.get(ctx, query, args...)
}

// get выполняет запрос к базе данных и возвращает пользователя.
// Используется внутренними методами GetByUUID, GetByLogin и GetByEmail.
func (r *Repository) get(ctx context.Context, query string, args ...any) (*model.User, error) {
	var repoUser repoModel.User

	err := r.pool.
		QueryRow(ctx, query, args...).
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

// buildGetUserQuery собирает SQL-запрос для получения пользователя с указанным условием.
func buildGetUserQuery(condition sq.Sqlizer) (string, []any, error) {
	selectUserBuilder := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select(
			"uuid",
			"login",
			"email",
			"password_hash",
			"notification_methods",
			"created_at",
			"updated_at",
		).
		From("users").
		Limit(1)

	return selectUserBuilder.Where(condition).ToSql()
}
