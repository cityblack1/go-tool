package random

import (
	"math/rand"
	"time"
)

func RandFloat64() float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()
}

func RandFloat32() float32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float32()
}

func RandPerm(n int) []int {
	rand.Seed(time.Now().UnixNano())
	return rand.Perm(n)
}

func RandIntN(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func RandInt63n(n int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(n)
}

func ShuffleFloat64Slice(items []float64) []float64 {
	if len(items) <= 1 {
		return items
	}

	ii := RandPerm(len(items))

	randItems := make([]float64, len(items))

	for i, randIdx := range ii {
		randItems[i] = items[randIdx]
	}
	return randItems
}

func ShuffleStringSlice(items []string) []string {
	if len(items) <= 1 {
		return items
	}

	ii := RandPerm(len(items))

	randItems := make([]string, len(items))

	for i, randIdx := range ii {
		randItems[i] = items[randIdx]
	}
	return randItems
}

func ShuffleInt64Slice(items []int64) []int64 {
	if len(items) <= 1 {
		return items
	}

	ii := RandPerm(len(items))

	randItems := make([]int64, len(items))

	for i, randIdx := range ii {
		randItems[i] = items[randIdx]
	}
	return randItems
}

func GetRandomAvgEmbedding(size int64) []float64 {
	res := make([]float64, size)
	for i := range res {
		res[i] = rand.Float64()*2 - 1 //生成的向量每个元素随机范围[-1,1)
	}
	return res
}
