package fst

import (
	"testing"
)

//func FuzzBuild(f *testing.F) {
//	testcases := []string{"Hello, world", " ", "!12345"}
//	for _, tc := range testcases {
//		f.Add(tc) // Use f.Add to provide a seed corpus
//	}
//	f.Fuzz(func(t *testing.T, orig []string) {
//
//	})
//}

func TestNewFst(t *testing.T) {
	fst := NewFst()
	fst.Set([]byte("abc"), 999)
	fst.Set([]byte("abcd"), 2)
	fst.Set([]byte("z"), 10)
	output, ok := fst.Search([]byte("abc"))
	assert(output, 999, ok, true)
	output, ok = fst.Search([]byte("abcd"))
	assert(output, 2, ok, true)
	output, ok = fst.Search([]byte("ab"))
	assert(ok, false)
	output, ok = fst.Search([]byte("za"))
	assert(ok, false)
	output, ok = fst.Search([]byte("z"))
	assert(output, 10, ok, true)
}

func TestNewFst2(t *testing.T) {
	fst := NewFst()
	fst.Set([]byte("abc"), 999)
	fst.Set([]byte("adc"), 2)
	fst.Set([]byte("z"), 2)
	//fst.debug()
	//fmt.Println(1)
}
