package v1

import (
	"github.com/radiophysiker/microservices-homework/iam/internal/service"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/user/v1"
)

// API представляет API слой для user service
type API struct {
	pb.UnimplementedUserServiceServer
	userService service.UserService
}

// NewAPI создает новый экземпляр API
func NewAPI(userService service.UserService) *API {
	return &API{
		userService: userService,
	}
}
