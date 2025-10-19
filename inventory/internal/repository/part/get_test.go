package part

import (
	"github.com/stretchr/testify/require"

	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
)

// TestGetPart проверяет метод GetPart репозитория
func (s *RepositoryTestSuite) TestGetPart() {
	tests := []struct {
		name    string
		uuid    string
		wantErr bool
		errType error
	}{
		{
			name:    "success_existing_part",
			uuid:    mainEngineV8UUID,
			wantErr: false,
		},
		{
			name:    "error_non_existent_part",
			uuid:    "99999999-9999-9999-9999-999999999999",
			wantErr: true,
			errType: model.ErrPartNotFound,
		},
		{
			name:    "error_invalid_uuid_format",
			uuid:    "invalid-uuid",
			wantErr: true,
			errType: model.ErrPartNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			part, err := s.repo.GetPart(s.ctx, tt.uuid)

			if tt.wantErr {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tt.errType)
				require.Nil(s.T(), part)
			} else {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), part)
				require.Equal(s.T(), tt.uuid, part.UUID)
			}
		})
	}
}
