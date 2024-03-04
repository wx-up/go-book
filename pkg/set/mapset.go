package set

type MapSet[T comparable] struct {
	m map[T]struct{}
}

func NewMapSet[T comparable](size int) *MapSet[T] {
	return &MapSet[T]{
		m: make(map[T]struct{}, size),
	}
}

func (m *MapSet[T]) Add(key ...T) {
	for _, k := range key {
		m.m[k] = struct{}{}
	}
}

func (m *MapSet[T]) Delete(key T) {
	delete(m.m, key)
}

func (m *MapSet[T]) Exist(key T) bool {
	_, ok := m.m[key]
	return ok
}

// Values 返回的元素顺序不固定
func (m *MapSet[T]) Values() []T {
	keys := make([]T, 0, len(m.m))
	for key := range m.m {
		keys = append(keys, key)
	}
	return keys
}
