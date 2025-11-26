package v1

import (
	"github.com/radiophysiker/microservices-homework/iam/internal/service"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/auth/v1"
)

// API представляет API слой для auth service
type API struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
}

// NewAPI создает новый экземпляр API
func NewAPI(authService service.AuthService) *API {
	return &API{
		authService: authService,
	}
}
