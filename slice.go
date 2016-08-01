package sys

func Contains(s []interface{}, v interface{}) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

func Merge(s1, s2 []interface{}) (c []interface{}) {
	c = append(s1, s2...)
	return
}

func SumInt(s []int) (sum int) {
	for _, v := range s {
		sum += v
	}
	return
}

func Int64Set(s []int64, sort bool) []int64 {
	size := len(s)
	if size == 0 || size == 1 {
		return s
	}
	m := make(map[int64]bool)
	for i := 0; i < size; i++ {
		m[s[i]] = true
	}
	rs := make([]int64, len(m))
	i := 0
	for key, _ := range m {
		rs[i] = key
		i++
	}
	if sort {
		size = len(rs)
		for j := 0; j < size; j++ {
			for k := 0; k < size; k++ {
				if rs[j] < rs[k] {
					rs[j], rs[k] = rs[k], rs[j]
				}
			}
		}
	}
	return rs
}
