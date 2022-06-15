package fst

import (
	"context"
	"sort"
)

type Dict interface {
	Search(word string) (int, bool)
	FuzzySearch(ctx context.Context, pattern string) <-chan string
}

type Node struct {
	next   map[byte]*Node
	id     int
	freeze bool
	val    byte
	output int
}

type Fst struct {
	DummyHead Node
	DummyTail Node
	Count     int
	countNode int
	autoIncr  int
}

type Pairs []Pair

type Pair struct {
	Word   string
	Output int
}

func (p Pairs) Len() int           { return len(p) }
func (p Pairs) Less(i, j int) bool { return p[i].Word < p[j].Word }
func (p Pairs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func NewPair(word string, output int) Pair {
	return Pair{Word: word, Output: output}
}

func NewFst(pairs ...Pair) *Fst {
	sort.Sort(Pairs(pairs))
	return &Fst{DummyHead: Node{}, DummyTail: Node{}, countNode: 0, Count: 0, autoIncr: 0}
}

func (f *Fst) Search(word string) (int, bool) {
	var ok bool
	curNode := &f.DummyHead
	sum := 0
	for _, char := range []byte(word) {
		if curNode, ok = curNode.next[char]; !ok {
			return 0, false
		}
		sum += curNode.output
	}
	if curNode != &f.DummyTail {
		return 0, false
	}
	return sum, true
}

func (f *Fst) FuzzySearch(ctx context.Context, pattern string) <-chan string {
	panic("implement me")
}

func (f *Fst) incr() int {
	f.autoIncr += 1
	return f.autoIncr
}
