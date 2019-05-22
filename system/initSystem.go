package system

import (
	//"database/sql"
	"fmt"
	"gin_blog/models"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pythonsite/iniConfig"
)

var DB *gorm.DB

type MySqlConfig struct {
	UserName string	`ini:"username"`
	Passwd string	`ini:"passwd"`
	DataBase string `ini:"database"`
	Host string	`ini:"host"`
	Port int	`ini:"port"`
}

type Page struct {
	PageSize int `ini:"page_size"`
}

type Config struct {
	MysqlConf  MySqlConfig `ini:"mysql"`
	Page Page `ini:"page"`
}

var ConfigConent *Config

func init() {
	ConfigConent = &Config{}
	err := iniConfig.UnmarshalFile("./conf/config.ini", ConfigConent)
	if err != nil {
		return
	}
	fmt.Printf("init config success! Config :%#v\n", ConfigConent)
	var ConsoleConfig=`{"level":7}`
	err = logs.SetLogger(logs.AdapterConsole,ConsoleConfig)
	if err != nil {
		logs.Error("init AdapterConsole logs error:%v", err)
		return
	}
	var FileConfig = `{"filename":"log/log.log","level":7}`
	err = logs.SetLogger(logs.AdapterFile, FileConfig)
	if err != nil {
		logs.Error("init AdapterFile logs error:%v", err)
		return
	}
	logs.Async(1e3)
	logs.Info("init logs Adapter success")
	mysqlDSN := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True",ConfigConent.MysqlConf.UserName, ConfigConent.MysqlConf.Passwd, ConfigConent.MysqlConf.Host, ConfigConent.MysqlConf.Port,ConfigConent.MysqlConf.DataBase)
	logs.Info("mysql dsn is:%v", mysqlDSN)
	db, err := gorm.Open("mysql", mysqlDSN)
	if err == nil {
		logs.Info("gorm open db success")
		DB = db
		db.AutoMigrate(&models.Page{}, &models.Post{}, &models.Tag{}, &models.User{}, &models.Comment{}, &models.Subscriber{}, &models.Link{})
		db.Model(&models.PostTag{}).AddUniqueIndex("uk_post_tag", "post_id", "tag_id")
	} else {
		logs.Error("gorm open mysql error:%v", err)
		return
	}

}



