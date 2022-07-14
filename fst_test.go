package fst

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

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

func TestFuzzySearch(t *testing.T) {
	fst := NewFst()
	fst.Set([]byte("aa"), 999)
	fst.Set([]byte("ab"), 2)
	fst.Set([]byte("abc"), 12390)
	fst.Set([]byte("b"), 2)
	fst.Set([]byte("dqwer"), 2)
	for kv := range fst.FuzzySearch(context.Background(), []byte(".b.")) {
		assert(kv.Word, []byte("abc"))
		assert(kv.Output, 12390)
	}
	count := 0
	for range fst.FuzzySearch(context.Background(), []byte("dqwe")) {
		count += 1
	}
	assert(count, 0)
}

func TestStop(t *testing.T) {
	fst := NewFst()
	fst.Set([]byte("a"), 999)
	fst.Set([]byte("b"), 2)
	fst.Set([]byte("c"), 12390)
	fst.Set([]byte("d"), 2)
	fst.Set([]byte("e"), 2)
	ctx, done := context.WithCancel(context.Background())
	defer done()
	count := 0
	for kv := range fst.FuzzySearch(ctx, []byte(".")) {
		fmt.Println(kv)
		done()
		count += 1
	}
	assert(count, 1)
}

func TestStop2(t *testing.T) {
	fst := NewFst()
	fst.Set([]byte("abc"), 999)
	fst.Set([]byte("abz"), 999)
	ctx, done := context.WithCancel(context.Background())
	defer done()
	count := 0
	for kv := range fst.FuzzySearch(ctx, []byte("ab.")) {
		fmt.Println(kv)
		done()
		count += 1
	}
	assert(count, 1)
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

func TestSetOutput6(t *testing.T) {
	fst := NewFst()
	fst.Set([]byte("mcvau"), 1)
	fst.Set([]byte("vjfv"), 2)
	fst.Set([]byte("vjfvau"), 3)
	fst.Set([]byte("zz"), 4)
	output, ok := fst.Search([]byte("vjfvau"))
	assert(ok, true)
	assert(output, 3)
}

func TestFreeze(t *testing.T) {
	fst := NewFst()
	fst.Set([]byte("aqwertyuiop"), 1)
	fst.Set([]byte("xqwertyuiop"), 1)
	fst.Set([]byte("zzz"), 1)

	fmt.Println(fst.Search([]byte("zzz")))
}

func TestEnsureSetWordKeepIncr(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {

		}
	}()
	fst := NewFst()
	fst.Set([]byte("z"), 3)
	fst.Set([]byte("a"), 1)
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
