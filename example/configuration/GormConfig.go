package configuration

import (
	"fmt"
	"github.com/happylusn/lithot-gin/lithot"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type GormConfig struct {
	SysConfig *lithot.SysConfig `inject:"-"`
}

func NewGormConfig() *GormConfig  {
	return &GormConfig{}
}

func (g *GormConfig) GormDB() *gorm.DB {
	dialect := getConfigValue(g.SysConfig.Config, []string{"db", "dialect"}, 0).(string)
	host := getConfigValue(g.SysConfig.Config, []string{"db", "host"}, 0)
	port := getConfigValue(g.SysConfig.Config, []string{"db", "port"}, 0)
	username := getConfigValue(g.SysConfig.Config, []string{"db", "username"}, 0)
	password := getConfigValue(g.SysConfig.Config, []string{"db", "password"}, 0)
	database := getConfigValue(g.SysConfig.Config, []string{"db", "database"}, 0)
	db, err := gorm.Open(dialect,
		fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local", username, password, host, port, database))
	if err != nil {
		log.Fatal(err)
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(5)                   //最大空闲数
	db.DB().SetMaxOpenConns(10)                  //最大打开连接数
	db.DB().SetConnMaxLifetime(time.Second * 30) //空闲连接生命周期
	return db
}

func getConfigValue(m lithot.UserConfig, prefix []string, index int) interface{} {
	res := lithot.GetConfigValue(m, prefix, index)
	if res == nil {
		panic("missing database configuration")
	}
	return res
}
