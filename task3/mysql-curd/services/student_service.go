package services

import (
	"fmt"
	"mysql/models"

	"gorm.io/gorm"
)

// Service 提供学生相关的CRUD操作
type Service struct {
	db *gorm.DB
}

// New 创建新的学生服务实例
func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreateBatch 批量插入学生记录
func (s *Service) CreateBatch(students []models.Student) error {
	result := s.db.Create(&students)
	if result.Error != nil {
		return fmt.Errorf("插入记录失败: %w", result.Error)
	}
	fmt.Printf("插入成功，插入数据行数为: %d\n", result.RowsAffected)
	return nil
}

// UpdateGradeByName 根据姓名更新学生年级
func (s *Service) UpdateGradeByName(name string, newGrade string) error {
	result := s.db.Model(&models.Student{}).Where("name = ?", name).Update("Grade", newGrade)
	if result.Error != nil {
		return fmt.Errorf("更新失败: %w", result.Error)
	}
	fmt.Printf("更新成功，受影响的记录数: %d\n", result.RowsAffected)
	return nil
}

// DeleteByAgeLessThan 根据年龄删除学生记录
func (s *Service) DeleteByAgeLessThan(age int) error {
	// 先查询符合条件的记录
	var studentsToDelete []models.Student
	s.db.Where("age < ?", age).Find(&studentsToDelete)

	fmt.Println("将要删除的学生记录:")
	for _, student := range studentsToDelete {
		fmt.Printf("  ID: %d 姓名: %s 年龄: %d 年级: %s\n",
			student.ID, student.Name, student.Age, student.Grade)
	}

	// 执行删除操作
	result := s.db.Where("age < ?", age).Delete(&models.Student{})
	if result.Error != nil {
		return fmt.Errorf("删除失败: %w", result.Error)
	}
	fmt.Printf("删除成功，受影响的记录数: %d\n", result.RowsAffected)
	return nil
}

// FindAll 查询所有学生记录
func (s *Service) FindAll() ([]models.Student, error) {
	var students []models.Student
	result := s.db.Find(&students)
	if result.Error != nil {
		return nil, fmt.Errorf("查询失败: %w", result.Error)
	}
	return students, nil
}

// FindByName 根据姓名查询学生
func (s *Service) FindByName(name string) (models.Student, error) {
	var student models.Student
	result := s.db.Where("name = ?", name).First(&student)
	if result.Error != nil {
		return models.Student{}, fmt.Errorf("查询失败: %w", result.Error)
	}
	return student, nil
}

// FindByAgeGreaterThan 根据年龄查询学生
func (s *Service) FindByAgeGreaterThan(age int) ([]models.Student, error) {
	var students []models.Student
	result := s.db.Where("age > ?", age).Find(&students)
	if result.Error != nil {
		return nil, fmt.Errorf("查询失败: %w", result.Error)
	}
	return students, nil
}
