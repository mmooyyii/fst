package fst

import (
	"testing"
)

func Test_smallByteMap_setEdge(t *testing.T) {

	bm := NewByteMap()
	bm.setEdge(0, &Edge{})
	bm.setEdge(8, &Edge{})
	bm.setEdge(2, &Edge{})
	bm.setEdge(23, &Edge{})
	ss := make([]byte, 4)
	for i, e := range bm.forloop() {
		ss[i] = e.byte
	}
	assert(ss, []byte{0, 2, 8, 23})
}
