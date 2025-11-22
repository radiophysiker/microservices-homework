package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/radiophysiker/microservices-homework/iam/internal/converter"
	"github.com/radiophysiker/microservices-homework/iam/internal/model"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/user/v1"
)

// GetUser обрабатывает запрос на получение информации о пользователе
func (a *API) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := a.userService.Get(ctx, req.UserUuid)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	protoUser := converter.ToProtoUser(user)

	return &pb.GetUserResponse{
		User: protoUser,
	}, nil
}
