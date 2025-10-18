package part

import (
	"errors"

	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
	"github.com/stretchr/testify/require"
)

// TestGetPart проверяет метод GetPart с различными сценариями
func (s *ServiceTestSuite) TestGetPart() {
	tests := []struct {
		name      string
		uuid      string
		setupMock func()
		wantPart  *model.Part
		wantErr   error
		checkErr  func(err error)
	}{
		{
			name: "success",
			uuid: "p1",
			setupMock: func() {
				want := &model.Part{
					UUID:        "p1",
					Name:        "Bolt",
					Description: "Description",
					Price:       10,
				}
				s.repo.EXPECT().
					GetPart(s.ctx, "p1").
					Return(want, nil).
					Once()
			},
			wantPart: &model.Part{
				UUID:  "p1",
				Name:  "Bolt",
				Price: 10,
			},
			wantErr: nil,
		},
		{
			name: "invalid_uuid_empty",
			uuid: "",
			setupMock: func() {
			},
			wantPart: nil,
			wantErr:  model.ErrInvalidUUID,
		},
		{
			name: "part_not_found",
			uuid: "nonexistent",
			setupMock: func() {
				s.repo.EXPECT().
					GetPart(s.ctx, "nonexistent").
					Return(nil, model.ErrPartNotFound).
					Once()
			},
			wantPart: nil,
			checkErr: func(err error) {
				require.Error(s.T(), err)
				require.Contains(s.T(), err.Error(), "failed to get part")
				require.True(s.T(), errors.Is(err, model.ErrPartNotFound))
			},
		},
		{
			name: "repository_error",
			uuid: "p3",
			setupMock: func() {
				repoErr := errors.New("database connection failed")
				s.repo.EXPECT().
					GetPart(s.ctx, "p3").
					Return(nil, repoErr).
					Once()
			},
			wantPart: nil,
			checkErr: func(err error) {
				require.Error(s.T(), err)
				require.Contains(s.T(), err.Error(), "failed to get part")
				require.Contains(s.T(), err.Error(), "database connection failed")
			},
		},
		{
			name: "invalid_uuid_with_special_chars",
			uuid: "p@#$%",
			setupMock: func() {
				s.repo.EXPECT().
					GetPart(s.ctx, "p@#$%").
					Return(nil, model.ErrPartNotFound).
					Maybe()
			},
			wantPart: nil,
			checkErr: func(err error) {
				require.Error(s.T(), err)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()
			part, err := s.service.GetPart(s.ctx, tt.uuid)

			if tt.checkErr != nil {
				tt.checkErr(err)
				require.Nil(s.T(), part)
			} else if tt.wantErr != nil {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tt.wantErr)
				require.Nil(s.T(), part)
			} else {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), part)
				require.Equal(s.T(), tt.wantPart.UUID, part.UUID)
				require.Equal(s.T(), tt.wantPart.Name, part.Name)
				require.Equal(s.T(), tt.wantPart.Price, part.Price)
				if tt.wantPart.Description != "" {
					require.Equal(s.T(), tt.wantPart.Description, part.Description)
				}
			}
		})
	}
}
