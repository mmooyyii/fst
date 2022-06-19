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

func increasingBytes(wordSize, sliceSize int) [][]byte {
	tmp := make(map[string]struct{})
	for i := 0; i < sliceSize; i++ {
		tmp[string(randomBytes(wordSize))] = struct{}{}
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

func randomBytes(length int) []byte {
	n := rand.Intn(length)
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
	b := increasingBytes(300, 300)
	for _, v := range b {
		fst.Set(v, rand.Intn(1000))
	}
}

func TestSetOutput(t *testing.T) {
	fst := NewFst()
	fst.Set([]byte("m"), 99)
	fst.Set([]byte("ma"), 100)
	fst.Set([]byte("mb"), 10)
	v, ok := fst.Search([]byte("m"))
	assert(99, v)
	assert(ok, true)
	v, ok = fst.Search([]byte("ma"))
	assert(100, v)
	assert(ok, true)
	v, ok = fst.Search([]byte("mb"))
	assert(10, v)
	assert(ok, true)
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

func TestSetOutput4(t *testing.T) {
	fst := NewFst()
	fst.Set([]byte("er"), 5199)
	fst.Set([]byte("hj"), 3901)
	fst.Set([]byte("hr"), 6310)
	fst.Set([]byte("o"), 1779)
	fmt.Println(fst.Search([]byte("hj")))
}

func TestSetOutput5(t *testing.T) {
	fst := NewFst()
	fst.Set([]byte("ecs"), 1399)
	fst.Set([]byte("es"), 879)
	fst.Set([]byte("is"), 7967)
	fst.Set([]byte("k"), 1965)
	fmt.Println(fst.Search([]byte("ecs")))
}

func TestFuzz(t *testing.T) {
	var seed int64
	var mock map[string]int
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("seed:", seed)
			fmt.Println(mock)
			panic(r)
		}
	}()
	for i := 0; i < 100; i++ {
		seed = time.Now().UnixMicro()
		//seed = int64(1655613820112461)
		rand.Seed(seed)
		fst := NewFst()
		mock = map[string]int{}
		b := increasingBytes(100, 100)
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
