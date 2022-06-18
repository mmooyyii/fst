package fst

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
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

func increasingBytes(n int) [][]byte {
	tmp := make(map[string]struct{})
	for i := 0; i < n; i++ {
		tmp[string(randomBytes())] = struct{}{}
	}
	ans := make([][]byte, 0)
	for k := range tmp {
		if len(k) > 0 {
			ans = append(ans, []byte(k))
		}
	}
	sort.Slice(ans, func(i, j int) bool {
		return bytes.Compare(ans[i], ans[j]) == -1
	})
	return ans
}

func randomBytes() []byte {
	n := rand.Intn(5)
	ans := make([]byte, n)
	for i := 0; i < n; i++ {
		ans[i] = byte(rand.Intn(20))
	}
	return ans
}

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
	fst.Set([]byte("aa"), 999)
	fst.Set([]byte("ab"), 2)
	fst.Set([]byte("abc"), 2)
	fst.Set([]byte("b"), 2)
	fst.Set([]byte("d"), 2)

	//fst.debug()
	////fmt.Println(1)
}

func TestBuild(t *testing.T) {
	rand.Seed(time.Now().Unix())
	fst := NewFst()
	b := increasingBytes(3000)
	fmt.Println(b)
	for _, v := range b {
		fst.Set(v, rand.Intn(1000))
	}

}
