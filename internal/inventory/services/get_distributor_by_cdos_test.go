package services

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/platform/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Unit_Services_Inventory_GetDistributorByCdos(t *testing.T) {

	type input struct {
		cdos GetDistributorByCdosDto
	}

	type result struct {
		getDistributorByCdos    inventory.Distributor
		getDistributorByCdosErr error
	}

	type want struct {
		distributor inventory.Distributor
		err         error
	}

	testCases := map[string]struct {
		input  input
		result result
		want   want
	}{
		"should return valid distributor": {
			input: input{
				cdos: GetDistributorByCdosDto{Cdos: "0140"},
			},
			result: result{
				getDistributorByCdos: inventory.Distributor{
					Name: "Test",
				},
				getDistributorByCdosErr: nil,
			},
			want: want{
				distributor: inventory.Distributor{
					Name: "Test",
				},
				err: nil,
			},
		},
		"should return err": {
			input: input{
				cdos: GetDistributorByCdosDto{Cdos: "0140"},
			},
			result: result{
				getDistributorByCdos:    inventory.Distributor{},
				getDistributorByCdosErr: errors.New("err"),
			},
			want: want{
				distributor: inventory.Distributor{},
				err:         errors.New("err"),
			},
		},
	}

	for name := range testCases {
		test := testCases[name]
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			inventoryRepository := new(mocks.RepositoryInventory)
			inventoryRepository.On("GetDistributorByCdos", ctx, test.input.cdos.Cdos).Return(test.result.getDistributorByCdos, test.result.getDistributorByCdosErr)

			svc := NewGetDistributorByCdosService(inventoryRepository)

			distributor, err := svc.Handler(ctx, test.input.cdos)

			assert.Equal(t, test.want.distributor, distributor)
			assert.Equal(t, test.want.err, err)
		})
	}
}
