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