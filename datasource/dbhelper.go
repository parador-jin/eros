package datasource

import (
	"Eros/conf"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"sync"
)

var dbLock sync.Mutex
var masterInstance *xorm.Engine

func InstanceDbMaster() *xorm.Engine {
	if masterInstance != nil {
		return masterInstance
	}
	dbLock.Lock()
	defer dbLock.Unlock()
	if masterInstance != nil {
		return masterInstance
	}
	return NewDbMaster()
}

func NewDbMaster() *xorm.Engine {
	sourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", conf.DbMaster.User, conf.DbMaster.Pwd,
		conf.DbMaster.Host, conf.DbMaster.Port, conf.DbMaster.Database)
	instance, err := xorm.NewEngine(conf.DriverName, sourceName)
	if err != nil {
		log.Fatal("dbHelper.NewDbMaster NewEngine error, error info is ", err)
	}
	instance.ShowSQL(true)
	instance.ShowExecTime(true)

	masterInstance = instance
	return instance
}
