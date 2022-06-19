package main

import (
	"context"
	"fmt"
	"github.com/mmooyyii/fst"
)

func main() {
	// 创建fst
	FST := fst.NewFst()
	// 插入, 需要保证插入字符串按字典顺递增， 否则会panic
	// 字符串不允许出现通配符 . (ascII: 46)
	FST.Set([]byte("abc"), 2)
	FST.Set([]byte("abz"), 9)

	// 查询
	fmt.Println(FST.Search([]byte("abz")))       // 9,true
	fmt.Println(FST.Search([]byte("not found"))) // 0,false

	// 模糊查询，模糊查询中，不能对fst进行其他操作，否则可能导致数据错乱
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for kv := range FST.FuzzySearch(ctx, []byte("a..")) {
		// . 代表这一位上可以是任意字符
		fmt.Println(string(kv.Word), kv.Output)
		// 调用cancel可以停止该channel
		cancel()
	}
}
