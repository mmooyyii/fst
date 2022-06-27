# fst

Finite State Transducer

### Installation

```shell
go get -u github.com/mmooyyii/fst
```

### example usage

```go
package main

import (
	"context"
	"fmt"
	"github.com/mmooyyii/fst"
)

func main() {
	// 创建fst
	FST := fst.NewFst()
	// 插入, 需要保证插入字符串按字典顺递增， 否则会panic
	// 字符串不允许出现通配符 . (ascII: 46)
	FST.Set([]byte("abc"), 2)
	FST.Set([]byte("abz"), 9)

	// 查询
	fmt.Println(FST.Search([]byte("abz")))       // 9,true
	fmt.Println(FST.Search([]byte("not found"))) // 0,false

	// 模糊查询，模糊查询中，不能对fst进行其他操作，否则可能导致数据错乱
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for kv := range FST.FuzzySearch(ctx, []byte("a..")) {
		// . 代表这一位上可以是任意字符
		fmt.Println(string(kv.Word), kv.Output)
		// 调用cancel可以停止该channel
		cancel()
	}
}
```

### benchmark

```
go test -bench=. -benchmem=true -run=none
```

```shell
goos: darwin
goarch: arm64
pkg: github.com/mmooyyii/fst
BenchmarkBuildFst-8            1        2122579750 ns/op        1216812328 B/op 24449745 allocs/op
BenchmarkBuildMap-8            9         117583028 ns/op        76083436 B/op     832588 allocs/op
BenchmarkSearchFst-8     1674333               735.0 ns/op            11 B/op          0 allocs/op
BenchmarkSearchMap-8     4845972               247.6 ns/op            11 B/op          0 allocs/op
PASS
ok      github.com/mmooyyii/fst 25.901s
```
