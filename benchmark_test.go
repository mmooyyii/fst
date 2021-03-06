package fst

import (
	"testing"
)

func BenchmarkBuildFst(b *testing.B) {
	w := increasingBytes(20, 1000000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fst := NewFst()
		for _, v := range w {
			fst.Set(v, 1)
		}
	}
}

func BenchmarkBuildMap(b *testing.B) {
	w := increasingBytes(20, 1000000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		M := map[string]int{}
		for _, v := range w {
			M[string(v)] = 1
		}
	}
}

func init() {

}

func BenchmarkSearchFst(b *testing.B) {
	var fff *Fst
	w := increasingBytes(20, 1000000)
	fff = NewFst()
	for _, v := range w {
		fff.Set(v, 1)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fff.Search(randomBytes(20))
	}
}

func BenchmarkSearchMap(b *testing.B) {
	w := increasingBytes(20, 1000000)
	M := make(map[string]int)
	for _, v := range w {
		M[string(v)] = 1
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := string(randomBytes(20))
		_ = M[k]
	}
}
