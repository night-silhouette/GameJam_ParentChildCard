package Util

import "math/rand/v2"

func RandomRange(min, max int) int {
	return rand.IntN(max-min+1) + min
}
