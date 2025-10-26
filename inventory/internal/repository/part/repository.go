package part

import "go.mongodb.org/mongo-driver/mongo"

// Repository реализует интерфейс PartRepository
type Repository struct {
	collection *mongo.Collection
}

// NewRepository создает новый экземпляр Repository
func NewRepository(collection *mongo.Collection) *Repository {
	return &Repository{collection: collection}
}
