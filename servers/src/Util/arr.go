package Util

// VerifyIncludes a包含b
func VerifyIncludes[t comparable](a, b []t) bool {
	if len(b) == 0 {
		return true
	}

	set := make(map[t]struct{}, len(a))
	for _, v := range a {
		set[v] = struct{}{}
	}
	for _, v := range b {
		if _, exists := set[v]; !exists {
			return false
		}
	}
	return true
}
