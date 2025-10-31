package part

import (
	"github.com/stretchr/testify/require"

	"github.com/radiophysiker/microservices-homework/inventory/internal/model"
)

// TestListParts проверяет метод ListParts репозитория
func (s *RepositoryTestSuite) TestListParts() {
	testData := GetTestParts()
	allParts := testData

	tests := []struct {
		name      string
		filter    *model.Filter
		wantCount int
		wantParts []*model.Part
	}{
		{
			name:      "success_nil_filter_returns_all",
			filter:    nil,
			wantCount: len(allParts),
			wantParts: allParts,
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
			wantParts: []*model.Part{allParts[0], allParts[1]},
		},
		{
			name: "success_filter_empty_uuids_returns_all",
			filter: &model.Filter{
				UUIDs: []string{},
			},
			wantCount: len(allParts),
			wantParts: allParts,
		},
		{
			name: "success_filter_nonexistent_uuid",
			filter: &model.Filter{
				UUIDs: []string{"nonexistent-uuid"},
			},
			wantCount: 0,
			wantParts: []*model.Part{},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.repo.On("ListParts", s.ctx, tt.filter).Return(tt.wantParts, nil).Once()

			parts, err := s.repo.ListParts(s.ctx, tt.filter)

			require.NoError(s.T(), err)
			require.Len(s.T(), parts, tt.wantCount)
			require.Equal(s.T(), tt.wantParts, parts)
		})
	}
}
