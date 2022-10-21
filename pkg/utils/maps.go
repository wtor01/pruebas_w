package utils

import (
	"strconv"
)

func ParseMapKey[T any](m map[string]string, key string, parse func(value string) (T, error)) (T, error) {
	val, ok := m[key]
	result, err := parse(val)

	if !ok {
		return result, err
	}

	return result, err
}

func ParseFloat64(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}

func Merge[K comparable, V any](old map[K]V, newM map[K]V) {
	for k, v := range newM {
		old[k] = v
	}
}

func SafeMerge[K comparable, V any](old map[K]V, newM map[K]V) map[K]V {
	if old == nil {
		old = make(map[K]V)
	}
	if newM == nil {
		newM = make(map[K]V)
	}
	for k, v := range newM {
		old[k] = v
	}
	return old
}

func MapToSlice[T any](items map[string]T) []T {
	arr := make([]T, 0, len(items))

	for _, v := range items {
		arr = append(arr, v)
	}

	return arr
}

func MapKeys[K comparable, V any](m map[K]V) []K {
	arr := make([]K, 0, len(m))
	for k, _ := range m {
		arr = append(arr, k)
	}

	return arr
}

func CopyMap[K comparable, V any](m map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		result[k] = v
	}

	return result
}
