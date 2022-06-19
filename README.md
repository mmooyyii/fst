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
BenchmarkBuildFst-8            1        3085338625 ns/op        1584858312 B/op 30032963 allocs/op
BenchmarkBuildMap-8            5         207924067 ns/op        114186686 B/op   1370124 allocs/op
BenchmarkSearchFst-8     1667916               721.2 ns/op            11 B/op          0 allocs/op
BenchmarkSearchMap-8     5168216               231.0 ns/op            11 B/op          0 allocs/op
```