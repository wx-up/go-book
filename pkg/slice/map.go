package slice

func Map[Src any, Dst any](values []Src, f func(idx int, val Src) Dst) []Dst {
	res := make([]Dst, len(values))
	for idx, val := range values {
		res[idx] = f(idx, val)
	}
	return res
}
