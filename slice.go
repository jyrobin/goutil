package goutil

func IndexOfString(arr []string, str string) int {
	for i := range arr {
		if arr[i] == str {
			return i
		}
	}
	return -1
}

func ContainsString(arr []string, str string) bool {
	return IndexOfString(arr, str) >= 0
}

func CharIn(ch rune, chs ...rune) bool {
	return charIn(ch, chs)
}
func charIn(ch rune, chs []rune) bool {
	for _, c := range chs {
		if ch == c {
			return true
		}
	}
	return false
}

func IndexOfChar(s string, chs ...rune) int {
	return indexOfChar(s, chs)
}
func indexOfChar(s string, chs []rune) int {
	for i, ch := range s {
		if charIn(ch, chs) {
			return i
		}
	}
	return -1
}

func CutLeft(s string, chs ...rune) string {
	if idx := indexOfChar(s, chs); idx >= 0 {
		return s[:idx]
	}
	return s
}

func CutRight(s string, chs ...rune) string {
	if idx := indexOfChar(s, chs); idx >= 0 {
		return s[idx+1:]
	}
	return s
}

func CutHalf(s string, chs ...rune) (string, string) {
	if idx := indexOfChar(s, chs); idx >= 0 {
		return s[:idx], s[idx+1:]
	}
	return s, ""
}
