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
		ans[i] = byte('a' + rand.Intn(20))
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

func TestSetOutput(t *testing.T) {
	fst := NewFst()
	fst.Set([]byte("m"), 99)
	fst.Set([]byte("ma"), 100)
	fst.Set([]byte("mb"), 10)

	fmt.Println(fst.Search([]byte("m")))
	fmt.Println(fst.Search([]byte("ma")))
	fmt.Println(fst.Search([]byte("mb")))
}

func TestSetOutput2(t *testing.T) {
	fst := NewFst()
	fst.Set([]byte("a"), 4)
	fst.Set([]byte("eb"), 3)
	fst.Set([]byte("ef"), 2)
	fst.Set([]byte("eq"), 1)
	v, ok := fst.Search([]byte("eb"))
	assert(3, v)
	assert(ok, true)
}

func TestSetOutput3(t *testing.T) {
	fst := NewFst()
	fst.Set([]byte("pk"), 3)
	fst.Set([]byte("pmj"), 2)
	fst.Set([]byte("pmp"), 1)
	fmt.Println(fst.Search([]byte("pk")))
}

func TestDuiPai(t *testing.T) {
	var seed int64
	var mock map[string]int
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("seed:", seed)
			fmt.Println(mock)
			panic(r)
		}
	}()
	for i := 0; i < 1000; i++ {
		seed = time.Now().UnixMicro()
		rand.Seed(int64(seed))
		fst := NewFst()
		mock = map[string]int{}
		b := increasingBytes(1000)
		for _, v := range b {
			output := rand.Intn(10000)
			fst.Set(v, output)
			mock[string(v)] = output
		}
		for _, v := range b {
			target := mock[string(v)]
			answer, ok := fst.Search(v)
			assert(ok, true)
			assert(target, answer)
		}
	}
}
