package util

import (
	"strings"
	"errors"
	"path"
)

type OrmInfo struct {
	DriverName string
	Addr    string
	Options map[string]string
}

func ExtractURL(s string) (*OrmInfo, error) {
	info := &OrmInfo{Options: make(map[string]string)}

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
	addr := strings.SplitN(s, "://", 2)
	if len(addr) == 2 {
		info.DriverName = addr[0]
		info.Addr = path.Clean(addr[1])
	} else {
		info.Addr = s
	}

	info.Addr = strings.Replace(info.Addr, "@", "@tcp(",1) + ")"
	return info, nil
}