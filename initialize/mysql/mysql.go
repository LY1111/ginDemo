package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tag_data_sync/config"
	"tag_data_sync/global"
	"tag_data_sync/initialize/logger"
	"time"
)

//初始化mysql 入口
func InitMySQL() error {
	var err error
	global.TagDB, err = getDssPool(config.Config.MysqlMaster.Host, config.Config.MysqlMaster.User, config.Config.MysqlMaster.Password, config.Config.MysqlMaster.Database, config.Config.MysqlMaster.Charset)
	if err != nil {
		return err
	}
	return nil
}

func getDssPool(host, user, password, database, charset string) (*gorm.DB, error) {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local", user, password, host, database, charset)

	logger.Info("连接mysql: ", dsn)

	pool, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db, err := pool.DB()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(config.Config.MysqlMaster.MaxIdleConn)
	db.SetMaxOpenConns(config.Config.MysqlMaster.MaxOpenConn)
	db.SetConnMaxLifetime(time.Duration(config.Config.MysqlMaster.MaxLifetime))

	return pool, nil
}

/**
 * 释放连接池
 */
func CloseMysqlPool() error {
	if global.TagDB != nil {
		db, err := global.TagDB.DB()
		if err != nil {
			return err
		}
		db.Close()
	}
	return nil
}
