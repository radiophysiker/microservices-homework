package payment

import (
	"os"
	"testing"

	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

func TestMain(m *testing.M) {
	logger.SetNopLogger()

	code := m.Run()
	os.Exit(code)
}
