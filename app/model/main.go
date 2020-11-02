package model

import (
	"back/app/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

/**
 * 新建数据库实例
 */
func New() {
	dbURL := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		config.Config.Database.User,
		config.Config.Database.Password,
		config.Config.Database.Address,
		config.Config.Database.Port,
		config.Config.Database.Dbname,
	)

	//println(dbURL)

	var err error
	db, err = gorm.Open("mysql", dbURL)


	if err != nil {
		temp := fmt.Sprintf("Database connect error: %v", err.Error())
		panic(temp)
	} else {
		db.SingularTable(true)
	}

	return
}

