package main

func IncludeString(vs []string, t string) bool {
	return indexOfString(vs, t) >= 0
}

func indexOfString(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}
