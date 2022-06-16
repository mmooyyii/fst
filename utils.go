package fst

func Map[T1 any, T2 any](f func(T1) T2, slice []T1) []T2 {
	ans := make([]T2, len(slice))
	for i, v := range slice {
		ans[i] = f(v)
	}
	return ans
}
