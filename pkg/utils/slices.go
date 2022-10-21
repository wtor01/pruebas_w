package utils

func MapSlice[K any, P any](arr []K, fn func(item K) P) []P {
	result := make([]P, 0, len(arr))

	for _, item := range arr {
		result = append(result, fn(item))
	}

	return result
}

func FindSlice[K any](arr []K, fn func(item K) bool) (*K, int) {
	for i, item := range arr {
		if fn(item) {
			return &item, i
		}
	}

	return nil, -1
}

type Set[T comparable] struct {
	m map[T]struct{}
}

func NewSet[T comparable](arr []T) Set[T] {
	m := make(map[T]struct{})

	for _, k := range arr {
		m[k] = struct{}{}
	}
	return Set[T]{
		m: m,
	}
}

func (s *Set[T]) Add(key T) {
	if _, ok := s.m[key]; !ok {
		s.m[key] = struct{}{}
	}
}

func (s *Set[T]) Has(key T) bool {
	_, ok := s.m[key]

	return ok
}

func (s *Set[T]) Slice() []T {
	result := make([]T, 0, len(s.m))

	for key, _ := range s.m {
		result = append(result, key)
	}

	return result
}

func InSlice[T comparable](value T, arr []T) bool {

	for _, item := range arr {
		if item == value {
			return true
		}
	}

	return false
}
