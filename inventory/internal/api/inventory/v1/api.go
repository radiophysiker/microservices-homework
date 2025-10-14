package v1

import (
	"github.com/radiophysiker/microservices-homework/inventory/internal/service"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
)

type API struct {
	pb.UnimplementedInventoryServiceServer
	partService service.PartService
}

// NewAPI creates new API.
func NewAPI(partService service.PartService) *API {
	return &API{
		partService: partService,
	}
}
