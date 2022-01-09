package text

import "strings"

// BetweenStrings returns a string between the starting and ending string or an empty string if none could be found
func BetweenStrings(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return
	}
	return str[s : s+e]
}
