package v1

import (
	"context"

	"github.com/radiophysiker/microservices-homework/inventory/internal/converter"
	pb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListParts возвращает список деталей с возможностью фильтрации
func (a *API) ListParts(ctx context.Context, req *pb.ListPartsRequest) (*pb.ListPartsResponse, error) {
	filter := converter.ToServiceFilter(req.Filter)

	parts, err := a.partService.ListParts(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list parts: %v", err)
	}

	protoParts := converter.ToProtoParts(parts)
	return &pb.ListPartsResponse{Parts: protoParts}, nil
}
