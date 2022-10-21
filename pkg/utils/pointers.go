package utils

func Float64(pointerFloat *float64) float64 {
	if pointerFloat != nil {
		return *pointerFloat
	}

	return 0.
}

func CreateFloat64(value float64) *float64 {
	return &value
}
