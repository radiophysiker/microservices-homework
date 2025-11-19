package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/radiophysiker/microservices-homework/iam/internal/model"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/auth/v1"
)

// Login обрабатывает запрос на вход пользователя
func (a *API) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	sessionUUID, err := a.authService.Login(ctx, req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidCredentials):
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	if sessionUUID == "" {
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &pb.LoginResponse{
		SessionUuid: sessionUUID,
	}, nil
}
