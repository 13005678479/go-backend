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
	db := database.InitDB(&models.Student{})

	// 创建学生服务实例
	svc := services.New(db)

	// 定义多条学生记录
	students := []models.Student{
		{Name: "张三", Age: 20, Grade: "三年级"},
		{Name: "李四", Age: 19, Grade: "三年级"},
		{Name: "王五", Age: 14, Grade: "四年级"},
	}

	// 使用CRUD方法
	if err := svc.CreateBatch(students); err != nil {
		log.Fatal("插入失败:", err)
	}

	if err := svc.UpdateGradeByName("张三", "四年级"); err != nil {
		log.Fatal("更新失败:", err)
	}

	if err := svc.DeleteByAgeLessThan(15); err != nil {
		log.Fatal("删除失败:", err)
	}

	// 查询示例
	allStudents, err := svc.FindAll()
	if err != nil {
		log.Fatal("查询失败:", err)
	}
	fmt.Printf("所有学生数量: %d\n", len(allStudents))

	student, err := svc.FindByName("李四")
	if err != nil {
		log.Fatal("查询失败:", err)
	}
	fmt.Printf("查询到的学生: %s %d %s\n", student.Name, student.Age, student.Grade)
}
