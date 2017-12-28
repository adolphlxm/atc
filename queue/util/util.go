package util

import (
	"errors"
	"strings"
)

type UrlInfo struct {
	Addr    string
	Options map[string]string
}

func ExtractURL(s string) (*UrlInfo, error) {
	info := &UrlInfo{Options: make(map[string]string)}

	if c := strings.Index(s, "?"); c != -1 {
		for _, pair := range strings.Split(s[c+1:], "&") {
			l := strings.SplitN(pair, "=", 2)
			if len(l) != 2 || l[0] == "" || l[1] == "" {
				return nil, errors.New("connection option must be key=value: " + pair)
			}
			info.Options[l[0]] = l[1]
		}
		s = s[:c]
	}
	info.Addr = s
	return info, nil
}
