package Util

import (
	"math/rand/v2"
)

func RandomRange(min, max int) int {
	return rand.IntN(max-min+1) + min
}

// GetRandomElements 从一个集合中随机抽取 $n$ 个不重复的元素
func GetRandomElements[T any](collection []T, n int) []T {
	count := len(collection)
	if count == 0 || n <= 0 {
		return []T{}
	}
	if n >= count {
		result := make([]T, count)
		copy(result, collection)
		return result
	}

	temp := make([]T, count)
	copy(temp, collection)

	for i := 0; i < n; i++ {
		j := i + rand.IntN(count-i)
		temp[i], temp[j] = temp[j], temp[i]
	}

	return temp[:n]
}
func GetRandomNormalInRange(n1, n2 float64, stdDevRatio float64) float64 {
	if n1 > n2 {
		n1, n2 = n2, n1
	}

	mean := (n1 + n2) / 2.0           // 均值取区间正中间
	stdDev := (n2 - n1) * stdDevRatio // 标准差与区间大小成正比

	// 使用 v2 的新 PCG 生成器（也可以直接用全局的 rand.NormFloat64()）
	rng := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))

	for {
		val := rng.NormFloat64()*stdDev + mean
		if val >= n1 && val <= n2 {
			return val
		}
	}
}
