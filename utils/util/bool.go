package util

func Mutex(bools ...bool) bool {
	count := 0
	for _, b := range bools {
		if b {
			count++
			if count > 1 {
				return false
			}
		}
	}
	return true
}
