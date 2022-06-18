package fst

import (
	"bytes"
	"context"
	"sync/atomic"
)

var incr int32

type Dict interface {
	GetPre() []byte
	Set(word []byte, output int) error
	Search(word []byte) (int, bool)
	FuzzySearch(ctx context.Context, pattern []byte) <-chan []byte
}

type Edge struct {
	node   *Node
	output int
	stop   bool // 是否是一个单词的结束
}

type Node struct {
	incr   int32
	next   map[byte]*Edge // 记录下一个节点的位置
	output int            // 只有在这是最后一个点的时候，才加上这里的output
}

func NewNode() *Node {
	return &Node{next: make(map[byte]*Edge, 0), incr: atomic.AddInt32(&incr, 1)}
}

type Fst struct {
	DummyHead Node             // 虚拟头节点
	DummyTail Node             // 虚拟尾节点
	Count     int              // 总单词数
	preWord   []byte           // 填加的单词必须单调递增， 所以记录上一个成功添加的字符串。
	TailHash  map[Hash][]*Node // 用于快速寻找最长的后缀
	Unfreeze  []*Node          // 未冻结的节点
}

func NewFst() *Fst {
	head := NewNode()
	tail := NewNode()
	return &Fst{
		DummyHead: *head,
		DummyTail: *tail,
		Unfreeze:  []*Node{head},
		preWord:   make([]byte, 0),
		TailHash:  make(map[Hash][]*Node, 0),
	}
}

func (f *Fst) Set(word []byte, output int) {
	if bytes.Compare(word, f.preWord) != 1 {
		panic("word must be increasing")
	}
	n := LongestPrefix(f.preWord, word) + 1
	output = f.PutOutput(n-1, output)
	f.freeze(n - 1)
	f.Unfreeze = f.Unfreeze[:n]
	preNode := f.Unfreeze[n-1]
	for i, char := range word[n-1:] {
		node := NewNode()
		preNode.next[char] = &Edge{node: node, output: output, stop: i+n == len(word)}
		preNode = node
		output = 0
		f.Unfreeze = append(f.Unfreeze, node)
	}
	f.preWord = word
}

func (f *Fst) GetPreWord() []byte { return f.preWord }

func (f *Fst) Search(word []byte) (int, bool) {
	return f.search(&f.DummyHead, word)
}

func (f *Fst) search(curNode *Node, word []byte) (int, bool) {
	sum := 0
	var stop bool
	for _, char := range word {
		if edge, ok := curNode.next[char]; !ok {
			return 0, false
		} else {
			curNode = edge.node
			sum += edge.output
			stop = edge.stop
		}
	}
	sum += curNode.output
	if !stop {
		return 0, false
	}
	return sum, true
}

func (f *Fst) FuzzySearch(ctx context.Context, pattern []byte) <-chan []byte {
	panic("implement me")
}

func (f *Fst) PutOutput(n int, output int) int {
	forwardOutput := 0
	for i, char := range f.preWord[:n] {
		v := f.Unfreeze[i]
		edge := v.next[char]
		if forwardOutput > 0 {
			edge.output += forwardOutput
		}
		if edge.output <= output {
			output -= edge.output
		} else {
			edge.output, forwardOutput = output, edge.output-output
			output = 0
		}
	}
	// 如果加不到边上， 就直接加到node上
	if forwardOutput > 0 {
		f.Unfreeze[len(f.Unfreeze)-1].output += forwardOutput
	}
	return output
}

func (f *Fst) freeze(n int) {
	suffixHash := SuffixHash(f.preWord[n:])
	for i, char := range f.preWord[n:] {
		if i == 0 {
			continue
		}
		hash := suffixHash[i]
		node := f.Unfreeze[n+i]
		if tail, ok := f.GetTail(hash, f.preWord[n+i:]); ok {
			node.next[char].node = tail
			return
		}
		if node.next[char].output == 0 {
			f.SetTailCache(hash, f.Unfreeze[n+i])
		}
	}
}

func (f *Fst) SetTailCache(hash Hash, node *Node) {
	f.TailHash[hash] = append(f.TailHash[hash], node)
}

func (f *Fst) GetTail(hash Hash, word []byte) (*Node, bool) {
	for _, node := range f.TailHash[hash] {
		if _, ok := f.search(node, word); ok {
			return node, true
		}
	}
	return nil, false
}
