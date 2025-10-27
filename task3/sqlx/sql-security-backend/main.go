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
	db := database.InitDB(&models.Employee{})

	svc := services.New(db)

	// db.Create(&models.Employee{ID: 1, Name: "张三", Department: "开发部", Salary: 100})
	// db.Create(&models.Employee{ID: 2, Name: "李四", Department: "技术部", Salary: 200})

	// 测试查询技术部员工
	techEmps, err := svc.GetTechEmployeesByDepartment("技术部")
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("技术部员工列表：")
		for _, e := range techEmps {
			fmt.Printf("ID: %d, 姓名: %s, 部门: %s, 工资: %d\n", e.ID, e.Name, e.Department, e.Salary)
		}
	}
	// 测试查询工资最高的员工
	topEmp, err := svc.GetTopSalaryEmployee(db)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("\n工资最高的员工：\nID: %d, 姓名: %s, 部门: %s, 工资: %d\n", topEmp.ID, topEmp.Name, topEmp.Department, topEmp.Salary)
	}

}
