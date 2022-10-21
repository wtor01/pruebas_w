package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"bitbucket.org/sercide/data-ingestion/internal/master_tables/geographic"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"context"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Unit_Service_Geographic_GetGeographicsService(t *testing.T) {

	type want struct {
		list  []geographic.GeographicZones
		count int
		err   error
	}

	tests := map[string]struct {
		input GetAllDto
		want  want
	}{
		"should call GetAllGeographicZones and return same values": {
			input: GetAllDto{
				Q:           "test",
				Limit:       10,
				Offset:      nil,
				Sort:        nil,
				CurrentUser: auth.User{},
			},
			want: want{
				list: []geographic.GeographicZones{
					{
						ID:          uuid.NewV4(),
						Code:        "ES",
						Description: "test",
						CreatedAt:   time.Time{},
						CreatedBy:   "test",
						UpdatedBy:   "test",
						UpdatedAt:   time.Time{},
					},
				},
				count: 1,
				err:   nil,
			},
		},
	}

	for testName, _ := range tests {
		testCase := tests[testName]

		t.Run(testName, func(t *testing.T) {
			repository := new(mocks.RepositoryGeographic)
			ctx := context.Background()

			repository.On("GetAllGeographicZones", ctx, geographic.Search{
				Q:           testCase.input.Q,
				Limit:       testCase.input.Limit,
				Offset:      testCase.input.Offset,
				Sort:        testCase.input.Sort,
				CurrentUser: testCase.input.CurrentUser,
			}).Return(testCase.want.list, testCase.want.count, testCase.want.err)

			svc := NewGetGeographicRepositoryService(repository)

			result, count, err := svc.Handler(ctx, testCase.input)

			assert.Equal(t, testCase.want.list, result)
			assert.Equal(t, testCase.want.count, count)
			assert.Equal(t, testCase.want.err, err)
		})
	}
}
