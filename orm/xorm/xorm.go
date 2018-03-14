package xorm

import (
	"fmt"
	"sync"

	"github.com/adolphlxm/atc/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"strconv"
	"github.com/adolphlxm/atc/orm/util"
)

type Orm struct {
	db       map[string]*xorm.EngineGroup
	logLevel core.LogLevel
	mu       sync.Mutex
}

func NewEngineGroup() orm.Orm {
	orm := &Orm{
		db: make(map[string]*xorm.EngineGroup),
	}
	return orm
}

func (this *Orm) Open(aliasName string, dataSourceName []string) error {
	var driverName string
	var dataSourceNameSlice []string
	var (
		maxIdleConns int
		maxOpenConns int
	)

	var (
		charset string
		db      string
	)
	for key, addr := range dataSourceName {
		dns, err := util.ExtractURL(addr)
		if err != nil {
			return err
		}

		_db := dns.Options["db"]
		_charset := dns.Options["charset"]

		// Master
		if key == 0 {
			driverName = dns.DriverName
			maxIdleConns, _ = strconv.Atoi(dns.Options["maxIdleConns"])
			maxOpenConns, _ = strconv.Atoi(dns.Options["maxOpenConns"])
		}

		if _charset != "" {
			charset = _charset
		}
		if _db != "" {
			db = _db
		}
		//dataSourceNameSlice = append(dataSourceNameSlice, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", host, db, charset))
		dataSourceNameSlice = append(dataSourceNameSlice, fmt.Sprintf("%s/%s?charset=%s", dns.Addr, db, charset))

	}

	engineGroup, err := xorm.NewEngineGroup(driverName, dataSourceNameSlice)
	if err != nil {
		return err
	}

	engineGroup.SetMaxIdleConns(maxIdleConns)
	engineGroup.SetMaxOpenConns(maxOpenConns)

	// The Ping () for the database connection test
	if err = engineGroup.Ping(); err != nil {
		return err
	}

	// Default.
	this.db[aliasName] = engineGroup
	return nil
}

func (this *Orm) SetLevel(aliasName string, level string) {
	//this.db[aliasName].ShowSQL(true)
	this.db[aliasName].Logger().ShowSQL(false)
	this.db[aliasName].Logger().SetLevel(core.LOG_OFF)
	if level == "" {
		level = "LOG_OFF"
	}

	switch level {
	case "LOG_UNKNOWN":
		this.logLevel = core.LOG_UNKNOWN
	case "LOG_OFF":
		this.logLevel = core.LOG_OFF
		this.db[aliasName].ShowSQL(false)
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
	this.db[aliasName].SetLogLevel(this.logLevel)
}

func (this *Orm) Ping(aliasName string) error {
	return this.db[aliasName].Ping()
}

// Ping tests if database is alive
func (this *Orm) Clone(aliasName string) error {
	slave := make([]*xorm.Engine, 0)
	master, _ := this.db[aliasName].Master().Clone()

	for _, slaveEngine := range this.db[aliasName].Slaves() {
		engine, _ := slaveEngine.Clone()
		slave = append(slave, engine)
	}

	var err error
	this.db[aliasName], err = xorm.NewEngineGroup(master, slave)
	return err
}

func (this *Orm) Use(aliasName string) *xorm.EngineGroup {
	return this.db[aliasName]
}

// Register
func init() {
	orm.Register("xorm", NewEngineGroup)
}
