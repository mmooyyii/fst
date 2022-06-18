//go:build debugMode

package fst

import (
	"encoding/json"
	"fmt"
)

type DebugEdge struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Char   byte   `json:"char"`
	Output int    `json:"output"`
}

func (f *Fst) debug() {
	ans := make([]*DebugEdge, 0)
	debugDfs(&f.DummyHead, &ans)
	data, err := json.Marshal(ans)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

func debugDfs(node *Node, ans *[]*DebugEdge) {
	for char, next := range node.next {
		*ans = append(*ans, &DebugEdge{
			From:   fmt.Sprintf("%d", node.incr),
			To:     fmt.Sprintf("%d", next.node.incr),
			Char:   char,
			Output: next.output,
		})
		debugDfs(next.node, ans)
	}
}
