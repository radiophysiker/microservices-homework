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

// Register обрабатывает запрос на регистрацию нового пользователя
func (a *API) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	registrationInfo := req.Info
	userInfo := &model.UserInfo{
		Login:               registrationInfo.Info.Login,
		Email:               registrationInfo.Info.Email,
		NotificationMethods: converter.FromProtoNotificationMethods(registrationInfo.Info.NotificationMethods),
	}

	password := registrationInfo.Password

	userUUID, err := a.userService.Register(ctx, userInfo, password)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	if userUUID == "" {
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &pb.RegisterResponse{
		UserUuid: userUUID,
	}, nil
}
