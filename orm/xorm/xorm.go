package xorm

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/adolphlxm/atc/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

type Orm struct {
	db       map[string]*xorm.Engine
	logLevel core.LogLevel
	mu       sync.Mutex
}

func NewXorm() orm.Orm {
	return &Orm{db: make(map[string]*xorm.Engine)}
}

func (this *Orm) Open(aliasName, config string) error {
	var cf map[string]string
	err := json.Unmarshal([]byte(config), &cf)
	if err != nil {
		return err
	}
	if cf["driver"] == "" {
		cf["driver"] = "mysql"
	}

	// New a db manager according to the parameter.
	engine, err := xorm.NewEngine(cf["driver"], fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", cf["user"], cf["password"], cf["host"], cf["dbname"]))
	if err != nil {
		return err
	}

	// The Ping () for the database connection test
	if err = engine.Ping(); err != nil {
		return err
	}

	switch cf["loglevel"] {
	case "LOG_UNKNOWN":
		this.logLevel = core.LOG_UNKNOWN
	case "LOG_OFF":
		this.logLevel = core.LOG_OFF
	case "LOG_ERR":
		this.logLevel = core.LOG_ERR
	case "LOG_WARNING":
		this.logLevel = core.LOG_WARNING
	case "LOG_INFO":
		this.logLevel = core.LOG_INFO
	case "LOG_DEBUG":
		this.logLevel = core.LOG_DEBUG
	default:
		this.logLevel = core.LOG_DEBUG
	}

	// Default.
	engine.Logger().SetLevel(core.LOG_OFF)

	this.db[aliasName] = engine
	return nil
}

func (this *Orm) SetMaxIdleConns(aliasName string, conns int) {
	this.db[aliasName].SetMaxIdleConns(conns)
}

func (this *Orm) SetMaxOpenConns(aliasName string, conns int) {
	this.db[aliasName].SetMaxOpenConns(conns)
}

func (this *Orm) Debug(aliasName string, show bool) {
	this.db[aliasName].ShowSQL(show)
	if show {
		this.db[aliasName].Logger().SetLevel(this.logLevel)
	} else {
		this.db[aliasName].Logger().SetLevel(core.LOG_OFF)
	}
}

func (this *Orm) Ping(aliasName string) error {
	return this.db[aliasName].Ping()
}

func (this *Orm) Clone(aliasName string) error {
	var err error
	this.mu.Lock()
	defer this.mu.Unlock()

	this.db[aliasName], err = this.db[aliasName].Clone()
	return err
}

func (this *Orm) Use(aliasName string) *xorm.Engine {
	return this.db[aliasName]
}

// Register
func init() {
	orm.Register("xorm", NewXorm)
}
