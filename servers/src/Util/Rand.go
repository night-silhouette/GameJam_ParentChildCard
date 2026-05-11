package Util

import (
	"math/rand/v2"
)

func RandomRange(min, max int) int {
	return rand.IntN(max-min+1) + min
}

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
