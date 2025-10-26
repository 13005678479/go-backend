package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Student struct {
	ID    int    `gorm:"primary_key;auto_increment"` // 主键且自增
	Name  string `gorm:"type:varchar(255)"`          // 学生姓名（字符串类型）
	Age   int    // 学生年龄（整数类型）
	Grade string `gorm:"type:varchar(50)"` // 学生年级（字符串类型）
}

// 可选：指定表名（如果结构体名与表名不一致时需要）
func (Student) TableName() string {
	return "biz_students"
}

func InitDB(dst ...interface{}) *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:password@2023@tcp(106.52.240.187:33306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(dst...)

	//单个插入
	// student := Student{
	// 	Name:  "张三",
	// 	Age:   20,
	// 	Grade: "三年级",
	// }

	// 定义多条学生记录
	student := []Student{
		{Name: "张三", Age: 20, Grade: "三年级"},
		{Name: "李四", Age: 19, Grade: "三年级"},
		{Name: "王五", Age: 21, Grade: "四年级"},
	}

	result := db.Create(&student)
	if result.Error != nil {
		panic("插入记录失败: " + result.Error.Error())
	}

	// 插入成功后，student.ID 会被自动赋值为自增的主键值
	println("插入成功，新记录ID为:", student)

	return db
}

func main() {
	InitDB(&Student{})
}
