package part

import (
	"errors"

	"github.com/stretchr/testify/require"

	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
)

// TestListParts проверяет метод ListParts с различными сценариями
func (s *ServiceTestSuite) TestListParts() {
	tests := []struct {
		name      string
		filter    *model.Filter
		setupMock func()
		wantParts []*model.Part
		wantErr   error
		checkErr  func(err error)
	}{
		{
			name:   "success_with_filter",
			filter: &model.Filter{UUIDs: []string{"p1", "p2"}},
			setupMock: func() {
				want := []*model.Part{
					{UUID: "p1", Name: "Bolt", Price: 10},
					{UUID: "p2", Name: "Screw", Price: 20},
				}
				s.repo.EXPECT().
					ListParts(s.ctx, &model.Filter{UUIDs: []string{"p1", "p2"}}).
					Return(want, nil).
					Once()
			},
			wantParts: []*model.Part{
				{UUID: "p1", Name: "Bolt", Price: 10},
				{UUID: "p2", Name: "Screw", Price: 20},
			},
			wantErr: nil,
		},
		{
			name:   "success_with_nil_filter",
			filter: nil,
			setupMock: func() {
				want := []*model.Part{
					{UUID: "p1", Name: "Bolt", Price: 10},
					{UUID: "p2", Name: "Screw", Price: 20},
				}
				s.repo.EXPECT().
					ListParts(s.ctx, (*model.Filter)(nil)).
					Return(want, nil).
					Once()
			},
			wantParts: []*model.Part{
				{UUID: "p1", Name: "Bolt", Price: 10},
				{UUID: "p2", Name: "Screw", Price: 20},
			},
			wantErr: nil,
		},
		{
			name:   "success_empty_result",
			filter: &model.Filter{UUIDs: []string{"nonexistent"}},
			setupMock: func() {
				s.repo.EXPECT().
					ListParts(s.ctx, &model.Filter{UUIDs: []string{"nonexistent"}}).
					Return([]*model.Part{}, nil).
					Once()
			},
			wantParts: []*model.Part{},
			wantErr:   nil,
		},
		{
			name:   "success_with_empty_filter",
			filter: &model.Filter{},
			setupMock: func() {
				want := []*model.Part{
					{UUID: "p1", Name: "Bolt", Price: 10},
					{UUID: "p2", Name: "Screw", Price: 20},
					{UUID: "p3", Name: "Nut", Price: 5},
				}
				s.repo.EXPECT().
					ListParts(s.ctx, &model.Filter{}).
					Return(want, nil).
					Once()
			},
			wantParts: []*model.Part{
				{UUID: "p1", Name: "Bolt", Price: 10},
				{UUID: "p2", Name: "Screw", Price: 20},
				{UUID: "p3", Name: "Nut", Price: 5},
			},
			wantErr: nil,
		},
		{
			name:   "repository_error",
			filter: &model.Filter{UUIDs: []string{"p1"}},
			setupMock: func() {
				repoErr := errors.New("database connection failed")
				s.repo.EXPECT().
					ListParts(s.ctx, &model.Filter{UUIDs: []string{"p1"}}).
					Return(nil, repoErr).
					Once()
			},
			wantParts: nil,
			checkErr: func(err error) {
				require.Error(s.T(), err)
				require.Contains(s.T(), err.Error(), "failed to list parts")
				require.Contains(s.T(), err.Error(), "database connection failed")
			},
		},
		{
			name:   "repository_returns_nil",
			filter: &model.Filter{UUIDs: []string{"p1"}},
			setupMock: func() {
				s.repo.EXPECT().
					ListParts(s.ctx, &model.Filter{UUIDs: []string{"p1"}}).
					Return(nil, nil).
					Once()
			},
			wantParts: nil,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			parts, err := s.service.ListParts(s.ctx, tt.filter)

			if tt.checkErr != nil {
				tt.checkErr(err)
				require.Nil(s.T(), parts)
			} else if tt.wantErr != nil {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tt.wantErr)
				require.Nil(s.T(), parts)
			} else {
				require.NoError(s.T(), err)

				if tt.wantParts == nil {
					require.Nil(s.T(), parts)
				} else {
					require.Equal(s.T(), len(tt.wantParts), len(parts))

					for i, want := range tt.wantParts {
						require.Equal(s.T(), want.UUID, parts[i].UUID)
						require.Equal(s.T(), want.Name, parts[i].Name)
						require.Equal(s.T(), want.Price, parts[i].Price)
					}
				}
			}
		})
	}
}
