package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/radiophysiker/microservices-homework/inventory/internal/converter"
	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
)

func (a *API) GetPart(ctx context.Context, req *pb.GetPartRequest) (*pb.GetPartResponse, error) {
	part, err := a.partService.GetPart(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrInvalidUUID) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid uuid: %s", req.Uuid)
		}

		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with uuid %s not found", req.Uuid)
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	protoPart := converter.ToProtoPart(part)

	return &pb.GetPartResponse{Part: protoPart}, nil
}
