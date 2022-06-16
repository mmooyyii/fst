package fst

import (
	"bytes"
	"context"
	"errors"
)

type Dict interface {
	Set(word []byte, output int) error
	Search(word []byte) (int, bool)
	FuzzySearch(ctx context.Context, pattern []byte) <-chan []byte
}

type Node struct {
	next   map[byte]*Node // 记录下一个节点的位置
	output int
	stop   bool
	val    byte
}

func NewNode(val byte) *Node {
	return &Node{next: make(map[byte]*Node, 0), val: val}
}

func (node *Node) link(node2 *Node) {
	node.next[node2.val] = node2
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
	head := NewNode(0)
	tail := NewNode(0)
	return &Fst{
		DummyHead: *head,
		DummyTail: *tail,
		Count:     0,
		Unfreeze:  make([]*Node, 0),
		preWord:   make([]byte, 0),
	}
}

func (f *Fst) Set(word []byte, output int) error {
	if bytes.Compare(word, f.preWord) != 1 {
		return errors.New("word must be increasing")
	}
	idx, output := f.Forward(word, output)
	toFreeze := make([]*Node, 0)
	copy(toFreeze, f.Unfreeze[idx:])
	f.Unfreeze = f.Unfreeze[:idx]
	// 把新加入的节点加入未冻结中
	for _, v := range word[idx:] {
		f.Unfreeze = append(f.Unfreeze, NewNode(v))
	}
	// 把新加入的节点连起来
	for i := idx; i < len(f.Unfreeze); i++ {
		f.Unfreeze[i-1].link(f.Unfreeze[i])
	}
	f.MergeFreeze(toFreeze)
	f.preWord = word
	return nil
}

func (f *Fst) Search(word []byte) (int, bool) {
	return f.search(&f.DummyHead, word)
}

func (f *Fst) search(curNode *Node, word []byte) (int, bool) {
	var ok bool
	sum := 0
	for _, char := range word {
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

func (f *Fst) FuzzySearch(ctx context.Context, pattern []byte) <-chan []byte {
	panic("implement me")
}

// Forward 在未冻结的节点中寻找最长公共前缀，并且调整output
// 返回最长公共前缀之后的那个节点的索引和修改后的output
func (f *Fst) Forward(word []byte, output int) (int, int) {
	// 在未冻结部分向前传递的output
	outputForward := 0
	for i, v := range word {
		if len(f.Unfreeze) == i {
			f.DummyTail.output += outputForward
			return i, output
		}
		if f.Unfreeze[i].val != v {
			f.Unfreeze[i].output += outputForward
			return i, output
		}
		if f.Unfreeze[i].output <= output {
			output -= f.Unfreeze[i].output
		} else {
			f.Unfreeze[i].output -= output
			outputForward = output
			output = 0
		}
	}
	return len(f.Unfreeze), output
}

// MergeFreeze 把一串节点冻结
// 先找出最长匹配后缀 然后把toFreeze连上去
func (f *Fst) MergeFreeze(toFreeze []*Node) {
	node1, node2 := f.longestSuffix(toFreeze)
	node1.link(node2)
}

func (f *Fst) longestSuffix(toFreeze []*Node) (*Node, *Node) {
	word := Map(func(node *Node) byte { return node.val }, toFreeze)
	suffixHash := SuffixHash(toFreeze)
	for idx, hash := range suffixHash {
		if nodes, ok := f.TailHash[hash]; ok {
			for _, node := range nodes {
				if _, ok2 := f.search(node, word[idx:]); ok2 {
					return toFreeze[idx], node
				}
			}
		}
	}
	return toFreeze[len(toFreeze)-1], &f.DummyTail
}

func SuffixHash(toFreeze []*Node) []Hash {
	hash := Hash(0)
	ans := make([]Hash, len(toFreeze))
	for i := len(toFreeze) - 1; i >= 0; i++ {
		hash = hash.Append(toFreeze[i].val)
		ans[i] = hash
	}
	return ans
}
