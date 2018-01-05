package queue_nats_streaming

import (
	"errors"
	"strconv"
	"time"

	"github.com/adolphlxm/atc/queue/util"
	"github.com/nats-io/go-nats-streaming"
)

type dialInfo struct {
	// Address holds the addresses for the server.
	ClusterID string
	ClientID string

	Options []stan.Option
}

func parseURL(url string) (*dialInfo, error) {
	opt, err := util.ExtractURL(url)
	if err != nil {
		return nil, err
	}

    var info = &dialInfo{Options: make([]stan.Option, 0)}

    if opt.Options["cluster"] == "" {
        return nil, errors.New("queue: cluster id is required")
    }
    info.ClusterID = opt.Options["cluster"]
    delete(opt.Options, "cluster")

    if opt.Options["client"] == "" {
        return nil, errors.New("client id is required")
    }
    info.ClientID = opt.Options["client"]
    delete(opt.Options, "client")

	for k, v := range opt.Options {
		switch k {
		case "maxPubAcksInflight":
			if maxRe, err := strconv.Atoi(v); err != nil {
				return nil, errors.New("bad value for maxPubAcksInflight: " + v)
			} else {
				info.Options = append(info.Options, stan.MaxPubAcksInflight(maxRe))
			}
		case "connectTimeout":
			if to, err := strconv.Atoi(v); err != nil {
				return nil, errors.New("bad value for connectTimeout: " + v)
			} else {
				info.Options = append(info.Options, stan.ConnectWait(time.Duration(to)*time.Millisecond))
			}
		case "ackTimeout":
			if to, err := strconv.Atoi(v); err != nil {
				return nil, errors.New("bad value for ackTimeout: " + v)
			} else {
				info.Options = append(info.Options, stan.PubAckWait(time.Duration(to)*time.Millisecond))
			}
		default:
			return nil, errors.New("unsupported connection URL option: " + k + "=" + v)
		}
	}
	info.Options = append(info.Options, stan.NatsURL(opt.Addr))
	return info, nil
}
