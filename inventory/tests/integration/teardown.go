//go:build integration

package integration

import (
	"context"
)

// Teardown выполняет очистку после всех тестов
func Teardown(ctx context.Context, env *TestEnvironment) error {
	// Очищаем тестовые данные
	if env.Collection != nil {
		if err := CleanupTestData(ctx, env.Collection); err != nil {
			return err
		}
	}

	// Останавливаем окружение
	return env.Teardown(ctx)
}
