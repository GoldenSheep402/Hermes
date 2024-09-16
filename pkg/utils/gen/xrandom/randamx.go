package xrandom

import (
	"math/rand"
)

type Type int

const (
	RandNum   Type = iota // 纯数字
	RandLower             // 小写字母
	RandUpper             // 大写字母
	RandAll               // 数字、大小写字母
)

// GetRandom 随机字符串 0 纯数字 1 小写字母 2 大写字母 3 数字、大小写字母 默认3
func GetRandom(size int, kind ...Type) string {
	if len(kind) == 0 {
		kind = append(kind, RandAll)
	}
	k := kind[0]
	iKind, kinds, result := k, [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
	isAll := k > 2 || k < 0
	for i := 0; i < size; i++ {
		if isAll { // random ikind
			iKind = Type(rand.Intn(3))
		}
		scope, base := kinds[iKind][0], kinds[iKind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}
