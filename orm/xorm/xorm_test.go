package xorm

import (
	"testing"
	"time"
	"github.com/lxmgo/atc/orm"
	"fmt"
)

type Test1 struct {
	Id int64 `xorm:"pk autoincr"`
	Number int64 `xorm:"int(11)"`
}

func _runEngine(t *testing.T) orm.Orm {
	xorm, _ := orm.NewOrm("xorm")
	err := xorm.Open("test_w",`{"driver":"mysql","host":"127.0.0.1:3306","user":"root","password":"123456","dbname":"test"}`)
	if err != nil {
		t.Errorf("Orm Open err:%v",err.Error())
	}

	return xorm
}

func TestConnect(t *testing.T) {
	xorm := _runEngine(t)
	xorm.Debug("test_w",true)
	engine := xorm.Use("test_w")

	if err := engine.Ping(); err != nil {
		t.Errorf("Orm ping failed err:%v",err.Error())
	}

	engine.Logger().Infof("Orm ping is success.")
}

func TestReconnect(t *testing.T) {
	xorm := _runEngine(t)
	xorm.Debug("test_w",true)
	engine := xorm.Use("test_w")

	engine.Logger().Infof("Please Start the database %v", engine.DriverName())
	time.Sleep(10 * time.Second)

	// Reconnect.
	if err := xorm.Clone("test_w"); err != nil {
		t.Errorf("Orm clone failed err:%v.",err.Error())
	}

	test1 := new(Test1)
	if _, err := engine.Where("id=?",1).Get(test1); err != nil {
		t.Errorf("Orm database failed err:%v",err.Error())
	}
	fmt.Println(test1)
	if err := engine.Ping(); err != nil {
		t.Errorf("Orm Reconnect ping failed err:%v",err.Error())
	}

	engine.Logger().Infof("Finish.")

}
