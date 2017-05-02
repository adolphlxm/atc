package orm

import (
	"fmt"

	"github.com/go-xorm/xorm"
)

type OrmFunc func() Orm

type Orm interface {
	// new a db manager according to the parameter. Currently support four for xorm
	// Param:
	//	1.
	//	2. {"driver":"mysql","host":"127.0.0.1:3306","user":"root","password":"123456","dbname":"test"}
	Open(aliasName, config string) error

	SetMaxIdleConns(aliasName string, conns int)

	// SetMaxOpenConns is only available for go 1.2
	SetMaxOpenConns(aliasName string, conns int)

	Debug(aliasName string, show bool)

	Ping(aliasName string) error

	Clone(aliasName string) error

	// Xorm
	Use(aliasName string) *xorm.Engine
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
