package orm

import (
	"fmt"

	"github.com/go-xorm/xorm"
)

type OrmFunc func() Orm

type Orm interface {
	// new a db manager according to the parameter. Currently support four for xorm
	Open(aliasName string, dataSourceName []string) error

	Ping(aliasName string) error

	Clone(aliasName string) error

	// Xorm
	Use(aliasName string) *xorm.EngineGroup
}

var adapters = make(map[string]OrmFunc)

func Register(name string, adapter OrmFunc) {
	if adapter == nil {
		panic("ATC orm: Register handler is nil")
	}
	if _, found := adapters[name]; found {
		panic("ATC orm: Register failed for handler " + name)
	}
	adapters[name] = adapter
}

func NewOrm(adapterName string) (Orm, error) {
	if handler, ok := adapters[adapterName]; ok {
		return handler(), nil
	} else {
		return nil, fmt.Errorf("ATC orm: unknown adapter name %s failed.", adapterName)
	}
}
