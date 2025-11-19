package utils

import "math/rand"

func RandomRangeInt64(start, end int64) int64 {
	diff := end - start
	num := rand.Int63n(diff)
	return start + num
}
