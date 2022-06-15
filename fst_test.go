package fst

import (
	"testing"
)

func TestNewFst(t *testing.T) {
	fst := NewFst(
		NewPair("mon", 2),
		NewPair("tues", 3),
		NewPair("thurs", 5))
	fst.debug()

}
