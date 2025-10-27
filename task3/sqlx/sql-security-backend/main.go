package main

import (
	"fmt"
	"log"
	"mysql/database"
	"mysql/models"
	"mysql/services"
)

func main() {
	// 初始化数据库连接
	db := database.InitDB(&models.Book{})

	// 示例：插入测试数据（可选，用于验证查询）
	testBooks := []models.Book{
		{Title: "Go 编程实战", Author: "张三", Price: 69.90},
		{Title: "Python 入门", Author: "李四", Price: 45.50},
		{Title: "Java 高级开发", Author: "王五", Price: 89.00},
	}
	db.Create(&testBooks) // 批量插入测试数据

	svc := services.New(db)

	// 执行查询：价格大于 50 元的书籍
	expensiveBooks, err := svc.GetBooksPriceGreaterThan50()
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("价格大于 50列表：")
		for _, e := range expensiveBooks {
			fmt.Printf("ID: %d, 姓名: %s, 部门: %s, 工资: %d\n", e.ID, e.Title, e.Author, e.Price)
		}
	}
}
