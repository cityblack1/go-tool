package partition


func Partition(src []int64, perSize int) [][]int64 {
	size := len(src) / perSize
	if len(src) % perSize > 0 {
		size ++
	}
	r := make([][]int64, size)
	start := 0
	for start < size {
		r[start] = make([]int64, 0)
		start ++
	}
	for i, i2 := range src {
		m := i%size
		r[m] = append(r[m], i2)
	}
	return r
}
