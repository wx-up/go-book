package set

type Set[T comparable] interface {
	Add(key ...T)
	Delete(key T)
	Exist(key T) bool
	Values() []T
}
