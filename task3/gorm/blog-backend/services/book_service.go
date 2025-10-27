package services

import (
	"fmt"
	"mysql/models"
	"mysql/utils"

	"gorm.io/gorm"
)

// Service 提供相关的查询操作
type Service struct {
	db *gorm.DB
}

// New 创建新的服务实例
func New(db *gorm.DB) *Service {
	utils.Info("创建新的服务实例")
	return &Service{db: db}
}

// 查询价格大于 50 元的书籍，返回 Book 结构体切片（类型安全）
func (s *Service) GetBooksPriceGreaterThan50() ([]models.Book, error) {
	utils.Info("开始查询价格大于50元的书籍")
	
	// 定义用于接收查询结果的 Book 结构体切片
	var books []models.Book
	result := s.db.Where("price > ?", 50).Find(&books)
	
	if result.Error != nil {
		utils.Error("查询价格大于50元的书籍失败: %v", result.Error)
		return nil, fmt.Errorf("查询失败: %w", result.Error)
	}
	
	utils.Info("查询成功，找到 %d 本价格大于50元的书籍", len(books))
	utils.Debug("查询结果: %+v", books)
	
	return books, nil
}
