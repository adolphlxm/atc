package xorm

import (
	"fmt"
	"github.com/adolphlxm/atc/orm"
	"testing"
)

type Test1 struct {
	Id int64 `xorm:"pk autoincr"`
	N1 int64 `xorm:"int(11)"`
}

func _runEngine(t *testing.T) orm.Orm {
	xorm, _ := orm.NewOrm("xorm")
	dataSourceNames := []string{"mysql://root:123456@127.0.0.1:3306//?charset=utf8&maxidleconns=1&maxopenconns=1&pingtime=30&db=test", "root:123456@?db=test2"}
	err := xorm.Open("t1", dataSourceNames)
	if err != nil {
		t.Errorf("Orm Open err:%v", err.Error())
	}

	return xorm
}

func TestConnect(t *testing.T) {
	xorm := _runEngine(t)
	engine := xorm.Use("t1")

	if err := engine.Ping(); err != nil {
		t.Errorf("Orm ping failed err:%v", err.Error())
	}

	engine.Logger().Infof("Orm ping is success.")
}

func TestReconnect(t *testing.T) {
	xorm := _runEngine(t)
	engine := xorm.Use("t1")

	engine.Logger().Infof("Please Start the database %v", engine.DriverName())
	//time.Sleep(10 * time.Second)

	// Reconnect.
	if err := xorm.Clone("t1"); err != nil {
		t.Errorf("Orm clone failed err:%v.", err.Error())
	}

	test1 := new(Test1)
	if _, err := engine.Where("id=?", 1).Get(test1); err != nil {
		t.Errorf("Orm database failed err:%v", err.Error())
	}
	fmt.Println(test1)
	if err := engine.Ping(); err != nil {
		t.Errorf("Orm Reconnect ping failed err:%v", err.Error())
	}

	engine.Logger().Infof("Finish.")

}
