# fst

Finite State Transducer

### Installation

```shell
go get -u github.com/mmooyyii/fst
```

### example usage

```go

```

### benchmark

```
go test bench=. benchmem=true -run=none
```

```shell
yimo@YideMacBook-Pro fst % go test -bench=. -benchmem=true -run=none
goos: darwin
goarch: arm64
pkg: github.com/mmooyyii/fst
BenchmarkBuildFst-8            1        3184254084 ns/op        1586697920 B/op 30075596 allocs/op
BenchmarkBuildMap-8            4         256846229 ns/op        123745708 B/op   1505547 allocs/op
PASS
ok      github.com/mmooyyii/fst 7.002s
```