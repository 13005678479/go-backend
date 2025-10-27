package services

import (
	"blog/internal/models"
	"blog/pkg/utils"

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

// 查询指定用户的所有文章及对应评论
func getUserPostsWithComments(db *gorm.DB, userID uint) (*models.User, error) {

}
