package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/radiophysiker/microservices-homework/iam/internal/converter"
	"github.com/radiophysiker/microservices-homework/iam/internal/model"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/auth/v1"
)

// Whoami обрабатывает запрос на получение информации о текущем пользователе
func (a *API) Whoami(ctx context.Context, req *pb.WhoamiRequest) (*pb.WhoamiResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	session, user, err := a.authService.Whoami(ctx, req.SessionUuid)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrSessionNotFound):
			return nil, status.Error(codes.NotFound, "session not found")
		case errors.Is(err, model.ErrInvalidSession):
			return nil, status.Error(codes.Unauthenticated, "invalid session")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	if session == nil || user == nil {
		return nil, status.Error(codes.Internal, "failed to get session or user")
	}

	protoSession := converter.ToProtoSession(session)
	protoUser := converter.ToProtoUser(user)

	return &pb.WhoamiResponse{
		Session: protoSession,
		User:    protoUser,
	}, nil
}
