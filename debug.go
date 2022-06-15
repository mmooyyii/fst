//go:build debugMode

package fst

import (
	"encoding/json"
	"fmt"
)

type DebugEdge struct {
	a      int
	b      int
	char   byte
	output int
}

func (f *Fst) debug() {
	ans := make([]DebugEdge, 0)
	debugDfs(&f.DummyHead, ans)
	data, err := json.Marshal(ans)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

func debugDfs(node *Node, ans []DebugEdge) {
	for char, next := range node.next {
		ans = append(ans, DebugEdge{
			a:      node.id,
			b:      next.id,
			char:   char,
			output: node.output,
		})
		debugDfs(next, ans)
	}

}
