package text

import (
	"fmt"
	"net/url"
	"strings"
)

// Nl2Br converts \n to HTML <br> tags
func Nl2Br(s string) string {
	return strings.Replace(s, "\n", "<br>", -1)
}

// LinkToHTMLLink converts regular links to HTML links
func LinkToHTMLLink(s string) string {
	s = strings.Replace(s, "\n", " \n ", -1)
	for _, p := range strings.Split(s, " ") {
		u, err := url.ParseRequestURI(p)
		if err == nil && u.Scheme != "" {
			s = strings.Replace(s, p, fmt.Sprintf("<a href=\"%s\">%s</a>", p, p), 1)
		}
	}
	return strings.Replace(s, " \n ", "\n", -1)
}
