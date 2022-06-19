package fst

import "testing"

func BenchmarkBuildFst(b *testing.B) {
	w := increasingBytes(20, 1000000)
	for i := 0; i < b.N; i++ {
		fst := NewFst()
		for _, v := range w {
			fst.Set(v, 1)
		}
	}
}

func BenchmarkBuildMap(b *testing.B) {
	w := increasingBytes(20, 1000000)
	for i := 0; i < b.N; i++ {
		M := map[string]int{}
		for _, v := range w {
			M[string(v)] = 1
		}
	}
}
