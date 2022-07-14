package fst

import (
	"bytes"
	"context"
)

type Dict interface {
	Set(word []byte, output int)
	Search(word []byte) (int, bool)
	FuzzySearch(ctx context.Context, pattern []byte) <-chan Kv
}

const WildCard = '.'

type Kv struct {
	Word   []byte
	Output int
}

type Edge struct {
	node   *Node
	output int
	stop   bool // 是否是一个单词的结束
}

type Node struct {
	next   map[byte]*Edge // 记录下一个节点的位置
	output int            // 只有在这是最后一个点的时候，才加上这里的output
}

func NewNode() *Node {
	return &Node{next: make(map[byte]*Edge, 0)}
}

type Fst struct {
	dummyHead       Node           // 虚拟头节点
	preWord         []byte         // 填加的单词必须单调递增， 所以记录上一个成功添加的字符串。
	tailHash        map[hash]*Node // 用于快速寻找最长的后缀
	unfreeze        []*Node        // 未冻结的节点
	PreNoOutputTail int            // 前一个字符串中, 没有任何output的后缀

}

func NewFst() *Fst {
	head := NewNode()
	return &Fst{
		dummyHead: *head,
		unfreeze:  []*Node{head},
		preWord:   make([]byte, 0),
		tailHash:  make(map[hash]*Node, 0),
	}
}

func (f *Fst) Set(word []byte, output int) {
	if bytes.Compare(word, f.preWord) != 1 {
		panic("word must be increasing")
	}
	if bytes.IndexByte(word, WildCard) != -1 {
		panic("wildcard is not allow")
	}
	n := longestPrefix(f.preWord, word) + 1
	output = f.PutOutput(n-1, output)
	f.freeze(n - 1)
	f.unfreeze = f.unfreeze[:n]
	preNode := f.unfreeze[n-1]
	f.PreNoOutputTail = n - 1

	for i, char := range word[n-1:] {
		node := NewNode()
		preNode.next[char] = &Edge{node: node, output: output, stop: i+n == len(word)}
		preNode = node
		output = 0
		f.unfreeze = append(f.unfreeze, node)
	}
	f.preWord = word
}

func (f *Fst) Search(word []byte) (int, bool) {
	return f.search(&f.dummyHead, word)
}

func (f *Fst) FuzzySearch(ctx context.Context, pattern []byte) <-chan Kv {
	ans := make(chan Kv, 0)
	go func() {
		f.fuzzySearch(ctx, &f.dummyHead, pattern, 0, []byte{}, 0, false, ans)
		close(ans)
	}()
	return ans
}

func (f *Fst) fuzzySearch(ctx context.Context, node *Node, pattern []byte, idx int, trace []byte, output int, stop bool, c chan Kv) bool {
	if idx == len(pattern) {
		if !stop {
			return true
		}
		ans := make([]byte, len(trace))
		copy(ans, trace)
		select {
		case <-ctx.Done():
			return false
		case c <- Kv{Word: ans, Output: output + node.output}:
			return true
		}
	}
	char := pattern[idx]
	if char == WildCard {
		for anyChar, edge := range node.next {
			trace = append(trace, anyChar)
			if !f.fuzzySearch(ctx, edge.node, pattern, idx+1, trace, output+edge.output, edge.stop, c) {
				return false
			}
			trace = trace[:len(trace)-1]
		}
	} else {
		if edge, ok := node.next[char]; !ok {
			return true
		} else {
			trace = append(trace, char)
			if !f.fuzzySearch(ctx, edge.node, pattern, idx+1, trace, output+edge.output, edge.stop, c) {
				return false
			}
			trace = trace[:len(trace)-1]
		}
	}
	return true
}

func (f *Fst) PutOutput(n int, output int) int {
	forwardOutput := 0
	for i, char := range f.preWord[:n] {
		v := f.unfreeze[i]
		edge := v.next[char]
		if forwardOutput > 0 {
			for _, e := range v.next {
				e.output += forwardOutput
			}
			forwardOutput = 0
		}
		if edge.output <= output {
			output -= edge.output
		} else {
			edge.output, forwardOutput = output, edge.output-output
			output = 0
		}
		if edge.stop {
			edge.node.output += forwardOutput
		}
	}
	if forwardOutput > 0 && n < len(f.preWord) {
		for _, e := range f.unfreeze[n].next {
			e.output += forwardOutput
		}
	}
	return output
}

func (f *Fst) freeze(n int) {
	if f.PreNoOutputTail > n {
		n = f.PreNoOutputTail
	}
	sh := suffixHash(f.preWord[n:])
	skipFirst := true
	for i, char := range f.preWord[n:] {
		hashValue := sh[i]
		node := f.unfreeze[n+i]
		if node.next[char].output+node.output != 0 {
			continue
		}
		if skipFirst {
			skipFirst = false
			continue
		}
		if tail, ok := f.getTail(hashValue, f.preWord[n+i:]); ok {
			f.unfreeze[n+i].next[char].node = tail.next[char].node
			return
		}
		if node.next[char].output+node.output == 0 {
			f.setTailCache(hashValue, node)
		}
	}
}

func (f *Fst) setTailCache(hash hash, node *Node) {
	f.tailHash[hash] = node
}

func (f *Fst) getTail(hash hash, word []byte) (*Node, bool) {
	if node, ok := f.tailHash[hash]; ok {
		if _, ok := f.search(node, word); ok {
			return node, true
		}
	}
	return nil, false
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
	if !stop {
		return 0, false
	}
	sum += curNode.output
	return sum, true
}
