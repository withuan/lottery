package datasource

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"lottery/conf"
	"sync"
)

var masterInstance *xorm.Engine

var dbLock sync.Mutex

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
	sourcename := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		conf.DbMaster.User,
		conf.DbMaster.Pwd,
		conf.DbMaster.Host,
		conf.DbMaster.Port,
		conf.DbMaster.Database)
	instance, err := xorm.NewEngine(conf.DriverName, sourcename)
	if err != nil {
		log.Fatal("dbhelper.NewDbMaster NewEngine error ", err)
		return nil
	}
	instance.ShowSQL(true)

	masterInstance = instance
	return instance
}
