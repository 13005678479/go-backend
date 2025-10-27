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

// 查询价格大于 50 元的书籍，返回 Book 结构体切片（类型安全）
func (s *Service) GetBooksPriceGreaterThan50() ([]models.Book, error) {
	// 定义用于接收查询结果的 Employee 结构体切片
	var books []models.Book
	result := s.db.Where("price > ?", 50).Find(&books)
	if result.Error != nil {
		return nil, fmt.Errorf("查询失败: %w", result.Error)
	}
	return books, nil
}
