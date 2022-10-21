package parsers

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"sort"
	"strconv"
)

func orderMeasuresClose(list []gross_measures.MeasureCloseWrite) {
	sort.SliceStable(list, func(i, j int) bool {
		ii, _ := strconv.Atoi(list[i].MeterSerialNumber)
		jj, _ := strconv.Atoi(list[j].MeterSerialNumber)
		return ii > jj
	})

	for index, _ := range list {
		a := list[index]
		sort.SliceStable(a.Periods, func(i, j int) bool {
			return a.Periods[i].Period > a.Periods[j].Period
		})
		list[index] = a
	}
}
