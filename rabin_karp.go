package fst

import (
	"math"
)

// 使用RabinKarp算法可在O(n)的时间复杂度下生成某字符串的每个后缀的哈希值, 用来辅助寻找最长后缀

const (
	prime = 16777619
	mod   = math.MaxUint32
)

type Hash uint32

func (h Hash) Append(char byte) Hash {
	// 此处char+1是为了防止char全为0的不同长度的字符串的hash值都相同
	return Hash((uint64(h)*prime + uint64(char) + 1) % mod)
}
