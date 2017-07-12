package pool

import (
	"testing"
)

func TestPool_Get(t *testing.T) {
	pool := &Pool{
		Dial: func() (interface{}, error) {
			return nil, nil
		},
		Close: func(c interface{}) error {

			return nil
		},
		MaxActive:   10,
		MaxIdle:     8,
		IdleTimeout: 1,
	}
	for i := 0; i < 12; i ++ {
		_, err := pool.Get()
		if err != nil {
			t.Error(err)
		}
		pool.Put(nil,false)
	}
}
