//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// InventoryTestSuite представляет основной тестовый suite для e2e тестов
type InventoryTestSuite struct {
	suite.Suite
	env    *TestEnvironment
	ctx    context.Context
	cancel context.CancelFunc
}

// SetupSuite выполняется один раз перед всеми тестами
func (s *InventoryTestSuite) SetupSuite() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 5*time.Minute)

	var err error
	s.env, err = Setup(s.ctx)
	s.Require().NoError(err, "Failed to setup test environment")

	// Вставляем тестовые данные
	testParts := GetTestParts()
	err = SetupTestData(s.ctx, s.env.Collection, testParts)
	s.Require().NoError(err, "Failed to setup test data")

	// Небольшая задержка, чтобы убедиться, что данные записаны в MongoDB
	time.Sleep(500 * time.Millisecond)
}

// TearDownSuite выполняется один раз после всех тестов
func (s *InventoryTestSuite) TearDownSuite() {
	if s.env != nil {
		// Создаем новый контекст для teardown, так как основной контекст может быть отменен
		teardownCtx, teardownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer teardownCancel()

		err := Teardown(teardownCtx, s.env)
		s.NoError(err, "Failed to teardown test environment")
	}

	if s.cancel != nil {
		s.cancel()
	}
}

// TestInventorySuite запускает все тесты в suite
func TestInventorySuite(t *testing.T) {
	suite.Run(t, new(InventoryTestSuite))
}
