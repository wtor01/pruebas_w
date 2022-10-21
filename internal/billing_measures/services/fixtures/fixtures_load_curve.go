package fixtures

import (
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/internal/process_measures"
)

//CCH-Complete-Valid
var processedLoadCurveComplete = []process_measures.ProcessedLoadCurve{
	//Period 1
	{
		Period:     measures.P1,
		Origin:     measures.STM,
		AI:         3333,
		AE:         3333,
		R1:         3333,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{
		Period:     measures.P1,
		Origin:     measures.STM,
		AI:         3333,
		AE:         3333,
		R1:         3333,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{
		Period:                   measures.P1,
		Origin:                   measures.STM,
		AI:                       3334,
		AE:                       3334,
		R1:                       3334,
		InvalidationCodes:        nil,
		ValidationStatusIsManual: false,
		Magnitudes:               []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},

	//Period 2
	{
		Period:     measures.P2,
		Origin:     measures.STM,
		AI:         4000,
		AE:         4000,
		R1:         4000,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{
		Period:     measures.P2,
		Origin:     measures.STM,
		AI:         4000,
		AE:         4000,
		R1:         4000,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{
		Period:     measures.P2,
		Origin:     measures.STM,
		AI:         4000,
		AE:         4000,
		R1:         4000,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},

	//Period 3

	{
		Period:     measures.P3,
		Origin:     measures.STM,
		AI:         1500,
		AE:         1500,
		R1:         1500,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{
		Period:     measures.P3,
		Origin:     measures.STM,
		AI:         1500,
		AE:         1500,
		R1:         1500,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{
		Period:     measures.P3,
		Origin:     measures.STM,
		AI:         1500,
		AE:         1500,
		R1:         1500,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{
		Period:     measures.P3,
		Origin:     measures.STM,
		AI:         1500,
		AE:         1500,
		R1:         1500,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
}

//CCH-Incomplete-InValid
var processedLoadCurveIncompleteInvalid = []process_measures.ProcessedLoadCurve{
	//Period 1
	{

		Period:     measures.P1,
		Origin:     measures.Filled,
		AI:         3333,
		AE:         3333,
		R1:         3333,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P1,
		Origin:     measures.Filled,
		AI:         3333,
		AE:         3333,
		R1:         3333,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P1,
		Origin:     measures.Filled,
		AI:         3334,
		AE:         3334,
		R1:         3334,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},

	//Period 2
	{

		Period:     measures.P2,
		Origin:     measures.Filled,
		AI:         4000,
		AE:         4000,
		R1:         4000,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P2,
		Origin:     measures.Filled,
		AI:         4000,
		AE:         4000,
		R1:         4000,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P2,
		Origin:     measures.Filled,
		AI:         4000,
		AE:         4000,
		R1:         4000,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},

	//Period 3

	{

		Period:     measures.P3,
		Origin:     measures.Filled,
		AI:         1500,
		AE:         1500,
		R1:         1500,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P3,
		Origin:     measures.Filled,
		AI:         1500,
		AE:         1500,
		R1:         1500,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P3,
		Origin:     measures.Filled,
		AI:         1500,
		AE:         1500,
		R1:         1500,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P3,
		Origin:     measures.Filled,
		AI:         1500,
		AE:         1500,
		R1:         1500,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
}

//CCH-Incomplete-Valid-Windows (P3 Windows)
var processedLoadCurveIncompleteWindows = []process_measures.ProcessedLoadCurve{
	//Period 1
	{

		Period:     measures.P1,
		Origin:     measures.Filled,
		AI:         3333,
		AE:         3333,
		R1:         3333,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P1,
		Origin:     measures.STM,
		AI:         3333,
		AE:         3333,
		R1:         3333,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P1,
		Origin:     measures.STM,
		AI:         3334,
		AE:         3334,
		R1:         3334,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},

	//Period 2
	{

		Period:     measures.P2,
		Origin:     measures.Filled,
		AI:         4000,
		AE:         4000,
		R1:         4000,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P2,
		Origin:     measures.STM,
		AI:         4000,
		AE:         4000,
		R1:         4000,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P2,
		Origin:     measures.STM,
		AI:         4000,
		AE:         4000,
		R1:         4000,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},

	//Period 3
	{

		Period:     measures.P3,
		Origin:     measures.STM,
		AI:         500,
		AE:         500,
		R1:         500,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P3,
		Origin:     measures.Filled,
		AI:         1000,
		AE:         1000,
		R1:         1000,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P3,
		Origin:     measures.Filled,
		AI:         1500,
		AE:         1500,
		R1:         1500,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P3,
		Origin:     measures.Filled,
		AI:         1500,
		AE:         1500,
		R1:         1500,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P3,
		Origin:     measures.Filled,
		AI:         1500,
		AE:         1500,
		R1:         1500,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
}

//CCH-Incomplete-Valid-NOWindows (P3 Windows)
var processedLoadCurveIncompleteNOWindows = []process_measures.ProcessedLoadCurve{
	//Period 1
	{

		Period:     measures.P1,
		Origin:     measures.Filled,
		AI:         3333,
		AE:         3333,
		R1:         3333,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P1,
		Origin:     measures.STM,
		AI:         3333,
		AE:         3333,
		R1:         3333,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P1,
		Origin:     measures.STM,
		AI:         3334,
		AE:         3334,
		R1:         3334,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},

	//Period 2
	{

		Period:     measures.P2,
		Origin:     measures.Filled,
		AI:         4000,
		AE:         4000,
		R1:         4000,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P2,
		Origin:     measures.STM,
		AI:         4000,
		AE:         4000,
		R1:         4000,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P2,
		Origin:     measures.STM,
		AI:         4000,
		AE:         4000,
		R1:         4000,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},

	//Period 3
	{

		Period:     measures.P3,
		Origin:     measures.STM,
		AI:         500,
		AE:         500,
		R1:         500,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P3,
		Origin:     measures.STM,
		AI:         1000,
		AE:         1000,
		R1:         1000,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P3,
		Origin:     measures.Filled,
		AI:         1500,
		AE:         1500,
		R1:         1500,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P3,
		Origin:     measures.Filled,
		AI:         1500,
		AE:         1500,
		R1:         1500,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
	{

		Period:     measures.P3,
		Origin:     measures.Filled,
		AI:         1500,
		AE:         1500,
		R1:         1500,
		Magnitudes: []measures.Magnitude{measures.AE, measures.AI, measures.R1},
	},
}
