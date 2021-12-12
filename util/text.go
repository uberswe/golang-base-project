// Package util provides utility functions for various things such as strings
package util

import (
	"fmt"
	"net/url"
	"strings"
)

func NL2BR(s string) string {
	return strings.Replace(s, "\n", "<br>", -1)
}

func StringLinkToHTMLLink(s string) string {
	s = strings.Replace(s, "\n", " \n ", -1)
	for _, p := range strings.Split(s, " ") {
		u, err := url.ParseRequestURI(p)
		if err == nil && u.Scheme != "" {
			s = strings.Replace(s, p, fmt.Sprintf("<a href=\"%s\">%s</a>", p, p), 1)
		}
	}
	return strings.Replace(s, " \n ", "\n", -1)
}

func GetStringBetweenStrings(str string, start string, end string) (result string) {
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
