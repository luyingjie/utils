package psql

import (
	"github.com/luyingjie/utils/os/cfg/conf"

	"github.com/jinzhu/gorm"

	//postgres数据驱动
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	//DB 数据库
	DB *gorm.DB
)

/*
OpenDB 开启数据库链接
*/
func OpenDB(driver, host, port, username, dbname, password, ssl, modeType string) error {
	var err error

	DB, err = gorm.Open(driver, "host="+host+" port="+port+" user="+username+" dbname="+dbname+" password="+password+" sslmode="+ssl)
	if err != nil {
		return err
	}
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
	DB.DB().SetConnMaxLifetime(3)
	DB.SingularTable(true)
	DB.LogMode(true)
	if modeType == "pro" {
		// 目前先去掉，后面可能不使用这个类
		// DB.SetLogger(logger.Logger.Engine)
	}
	return nil
}

/*
InitDB 初始化数据库连接
*/
func InitDB() error {
	confs := conf.GetByKeyString("psql")
	//redis.Open(confing.Redis.Host, confing.Redis.Port, confing.Redis.Password)
	err := OpenDB(
		confs["driver"],
		confs["host"],
		confs["port"],
		confs["user"],
		confs["dbname"],
		confs["password"],
		confs["ssl"],
		"dev",
	)
	//远程链接数据库的情况建议建表后注释掉，否则检查表会可能慢
	//DB.AutoMigrate(App{}, Account{}, AuthSource{})
	if err != nil {
		return err
	}

	return nil
}
