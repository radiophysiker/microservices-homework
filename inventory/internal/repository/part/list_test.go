package part

import (
	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
	"github.com/stretchr/testify/require"
)

// TestListParts проверяет метод ListParts репозитория
func (s *RepositoryTestSuite) TestListParts() {
	lenTaskList := len(testParts)
	tests := []struct {
		name      string
		filter    *model.Filter
		wantCount int
		check     func(parts []*model.Part)
	}{
		{
			name:      "success_nil_filter_returns_all",
			filter:    nil,
			wantCount: lenTaskList,
			check: func(parts []*model.Part) {
				require.NotNil(s.T(), parts)
				require.Len(s.T(), parts, lenTaskList)
			},
		},
		{
			name: "success_filter_by_specific_uuids",
			filter: &model.Filter{
				UUIDs: []string{
					mainEngineV8UUID,
					fuelTankUUID,
				},
			},
			wantCount: 2,
			check: func(parts []*model.Part) {
				require.NotNil(s.T(), parts)
				require.Len(s.T(), parts, 2)
			},
		},
		{
			name: "success_filter_empty_uuids_returns_all",
			filter: &model.Filter{
				UUIDs: []string{},
			},
			wantCount: lenTaskList,
			check: func(parts []*model.Part) {
				require.NotNil(s.T(), parts)
				require.Len(s.T(), parts, lenTaskList)
			},
		},
		{
			name: "success_filter_nonexistent_uuid",
			filter: &model.Filter{
				UUIDs: []string{"nonexistent-uuid"},
			},
			wantCount: 0,
			check: func(parts []*model.Part) {
				require.NotNil(s.T(), parts)
				require.Empty(s.T(), parts)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			parts, err := s.repo.ListParts(s.ctx, tt.filter)

			require.NoError(s.T(), err)
			require.Len(s.T(), parts, tt.wantCount)
			tt.check(parts)
		})
	}
}
