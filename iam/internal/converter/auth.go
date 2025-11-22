package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/radiophysiker/microservices-homework/iam/internal/model"
	commonpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/common/v1"
)

// ToProtoSession преобразует доменную модель Session в protobuf Session
func ToProtoSession(s *model.Session) *commonpb.Session {
	if s == nil {
		return nil
	}

	return &commonpb.Session{
		Uuid:      s.UUID,
		CreatedAt: timestamppb.New(s.CreatedAt),
		UpdatedAt: timestamppb.New(s.UpdatedAt),
		ExpiresAt: timestamppb.New(s.ExpiresAt),
	}
}
