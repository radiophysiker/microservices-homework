package integration

const (
	// TestNetworkName имя тестовой Docker сети
	TestNetworkName = "inventory-test-network"

	// TestMongoContainerName имя контейнера MongoDB для тестов
	TestMongoContainerName = "inventory-test-mongo"

	// TestMongoDatabase имя тестовой базы данных
	TestMongoDatabase = "inventory_test"

	// TestMongoUsername имя пользователя MongoDB для тестов
	TestMongoUsername = "testuser"

	// TestMongoPassword пароль MongoDB для тестов
	TestMongoPassword = "testpass"

	// TestMongoAuthDB база данных для аутентификации
	TestMongoAuthDB = "admin"

	// TestAppContainerName имя контейнера приложения для тестов
	TestAppContainerName = "inventory-test-app"

	// TestAppPort порт gRPC сервера приложения
	TestAppPort = "50051"

	// TestCollectionName имя коллекции для деталей
	TestCollectionName = "parts"
)
