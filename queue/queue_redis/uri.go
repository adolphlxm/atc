package queue_redis

import (
	"errors"
	"strconv"
	"time"

    "github.com/adolphlxm/atc/queue/util"
)

type dialInfo struct {
	// Address holds the addresses for the server.
	Url string

	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration

	IdleTimeout time.Duration

	MaxIdle   int
	MaxActive int
}

func parseURL(url string) (*dialInfo, error) {
	opt, err := util.ExtractURL(url)

	if err != nil {
		return nil, err
	}

	var (
		maxIdle, maxActive, connectTimeout, readTimeout, writeTimeout, idleTimeout int
	)
	for k, v := range opt.Options {
		switch k {
		case "maxIdle":
			if maxIdle, err = strconv.Atoi(v); err != nil {
				return nil, errors.New("bad value for maxIdle: " + v)
			}
		case "maxActive":
			if maxActive, err = strconv.Atoi(v); err != nil {
				return nil, errors.New("bad value for maxActive: " + v)
			}
		case "connectTimeout":
			if connectTimeout, err = strconv.Atoi(v); err != nil {
				return nil, errors.New("bad value for connectTimeout: " + v)
			}
		case "readTimeout":
			if readTimeout, err = strconv.Atoi(v); err != nil {
				return nil, errors.New("bad value for readTimeout: " + v)
			}
		case "writeTimeout":
			if writeTimeout, err = strconv.Atoi(v); err != nil {
				return nil, errors.New("bad value for writeTimeout: " + v)
			}
		case "idleTimeout":
			if idleTimeout, err = strconv.Atoi(v); err != nil {
				return nil, errors.New("bad value for idleTimeout: " + v)
			}
		default:
			return nil, errors.New("unsupported connection URL option: " + k + "=" + v)
		}
	}

	if _, ok := opt.Options["maxIdle"]; !ok {
		maxIdle = 1
	}

	if _, ok := opt.Options["maxActive"]; !ok {
		maxActive = 1
	}
	if _, ok := opt.Options["connectTimeout"]; !ok {
	    connectTimeout = 3000
    }

	info := dialInfo{
		Url:            opt.Addr,
		MaxIdle:        maxIdle,
		MaxActive:      maxActive,
		ConnectTimeout: time.Duration(connectTimeout) * time.Millisecond,
		ReadTimeout:    time.Duration(readTimeout) * time.Millisecond,
		WriteTimeout:   time.Duration(writeTimeout) * time.Millisecond,
		IdleTimeout:    time.Duration(idleTimeout) * time.Millisecond,
	}
	return &info, nil
}
