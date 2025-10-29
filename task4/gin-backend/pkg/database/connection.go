package database

import (
	"blog/pkg/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(dst ...interface{}) *gorm.DB {
	utils.Info("正在初始化数据库连接...")

	dsn := "root:password@2023@tcp(106.52.240.187:33306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	utils.Debug("数据库连接字符串: %s", dsn)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		utils.Error("数据库连接失败: %v", err)
		panic(err)
	}

	utils.Info("数据库连接成功")

	// 执行数据库迁移
	utils.Info("开始执行数据库迁移...")
	db.AutoMigrate(dst...)
	utils.Info("数据库迁移完成")

	return db
}
