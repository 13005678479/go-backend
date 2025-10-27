package services

import (
	"fmt"
	"mysql/models"

	"gorm.io/gorm"
)

// Service 提供相关的查询操作
type Service struct {
	db *gorm.DB
}

// New 创建新的服务实例
func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

// 查询所有技术部员工（映射到 []Employee）
func (s *Service) GetTechEmployeesByDepartment(department string) ([]models.Employee, error) {
	// 定义用于接收查询结果的 Employee 结构体切片
	var employees []models.Employee
	result := s.db.Where("department = ?", department).Find(&employees)
	if result.Error != nil {
		return nil, fmt.Errorf("查询失败: %w", result.Error)
	}
	return employees, nil
}

// 需求2：查询工资最高的员工（映射到 Employee）
func (s *Service) GetTopSalaryEmployee(db *gorm.DB) (models.Employee, error) {
	// 使用 GORM 的查询方法
	var emp models.Employee
	// 使用 Order 按工资降序排序，First 取第一条记录
	if err := db.Order("salary DESC").First(&emp).Error; err != nil {
		// 若查询出错，返回空 Employee 结构体和包含错误信息的 error
		return models.Employee{}, fmt.Errorf("查询最高工资员工失败: %v", err)
	}
	// 查询成功，返回员工结构体和 nil 错误
	return emp, nil
}
