package effective

func EffectiveString(strs ...string) bool {
	for _, str := range strs {
		if str == "" {
			return false
		}
	}
	return true
}
