//go:build integration

package integration

import (
	"context"
)

// Teardown выполняет очистку после всех тестов
func Teardown(ctx context.Context, env *TestEnvironment) error {
	if env.Collection != nil {
		if err := CleanupTestData(ctx, env.Collection); err != nil {
			return err
		}
	}

	return env.Teardown(ctx)
}


