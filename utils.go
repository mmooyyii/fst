package fst

import (
	"fmt"
	"math"
	"reflect"
)

// 使用RabinKarp算法可在O(n)的时间复杂度下生成某字符串的每个后缀的哈希值, 用来辅助寻找最长后缀

const (
	prime = 16777619
	mod   = math.MaxUint32
)

type hash uint32

func (h hash) append(char byte) hash {
	// 此处char+1是为了防止char全为0的不同长度的字符串的hash值都相同
	return hash((uint64(h)*prime + uint64(char) + 1) % mod)
}

type number interface{ int | int32 | int64 | uint }

func min[T number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func longestPrefix(a, b []byte) int {
	n := min(len(a), len(b))
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return n
}

// suffixHash 对于 "abc", 返回 [hash("abc"),hash("bc"),hash("c")]
func suffixHash(a []byte) []hash {
	ans := make([]hash, len(a))
	cur := hash(0)
	for i := len(a) - 1; i >= 0; i-- {
		cur = cur.append(a[i])
		ans[i] = cur
	}
	return ans
}

func assert(a ...interface{}) {
	for i := 0; i < len(a); i += 2 {
		if !reflect.DeepEqual(a[i], a[i+1]) {
			panic(fmt.Sprintf("ERROR: %v %v", a[i], a[i+1]))
		}
	}
}
