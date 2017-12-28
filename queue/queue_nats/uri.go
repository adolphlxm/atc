package queue_nats

import (
	"errors"
	"strconv"
	"time"

    "github.com/adolphlxm/atc/queue/util"
    "github.com/nats-io/go-nats"
)

type dialInfo struct {
	// Address holds the addresses for the server.
	Url string

	ConnectTimeout time.Duration

	Options []nats.Option
}

func parseURL(url string) (*dialInfo, error) {
	opt, err := util.ExtractURL(url)
	if err != nil {
		return nil, err
	}

	var (
		info = &dialInfo{Options:make([]nats.Option, 0)}
	)

	for k, v := range opt.Options {
		switch k {
		case "allowReconnect":
		    switch v {
            case "false","no","0":
                info.Options = append(info.Options, nats.NoReconnect())
            }
		case "name":
			if v != "" {
			    info.Options = append(info.Options, nats.Name(v))
            }
		case "noRandomize":
            switch v {
            case "true","yes","1":
                info.Options = append(info.Options, nats.DontRandomize())
            }
		case "maxReconnect":
			if maxRe, err := strconv.Atoi(v); err != nil {
				return nil, errors.New("bad value for maxReconnect: " + v)
			} else {
                info.Options = append(info.Options, nats.MaxReconnects(maxRe))
            }
        case "connectTimeout":
            if to, err := strconv.Atoi(v); err != nil {
                return nil, errors.New("bad value for connectTimeout: " + v)
            } else {
                info.Options = append(info.Options, nats.Timeout(time.Duration(to) * time.Millisecond))
            }
		default:
			return nil, errors.New("unsupported connection URL option: " + k + "=" + v)
		}
	}

    info.Url = opt.Addr
	return info, nil
}
