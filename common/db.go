package common

import "github.com/jinzhu/gorm"
import _ "github.com/jinzhu/gorm/dialects/sqlite"

var DB *gorm.DB

func newDb(dbPath string) *gorm.DB {
	db, err := gorm.Open("sqlite3", "test.db")

	//db.LogMode(true)

	if err != nil {
		panic("连接数据库失败:" + err.Error())
	}

	return db
}

func dbInit() {
	DB = newDb("./store.db")
}
