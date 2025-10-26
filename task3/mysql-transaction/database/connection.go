package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(dst ...interface{}) *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:password@2023@tcp(106.52.240.187:33306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(dst...)
	return db
}
