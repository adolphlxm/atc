package memcache

import (
	"testing"
	"github.com/adolphlxm/atc/cache"
	"time"
)

func TestNewMemCache(t *testing.T) {
	m, err := cache.NewCache("memcache","127.0.0.1:11211")
	if err != nil {
		t.Errorf("Memcache init err:%v",err.Error())
	}

	timeoutDuration := 10 * time.Second
	if err = m.Put("atc","atc framework",timeoutDuration); err != nil {
		t.Errorf("Memcache put err:%v",err.Error())
	}

	v, err := m.Get("atc")
	if err != nil {
		t.Errorf("Memcache Get err:%v",err.Error())
	}

	if string(v) != "atc framework" {
		t.Errorf("Memcache Get value error.")
	}

	if err := m.Delete("atc"); err != nil {
		t.Errorf("Memcache Delete err:%v",err.Error())
	}
}