package services

import (
	"blogV2/internal/models"

	"gorm.io/gorm"
)

// CommentService 评论服务
type CommentService struct {
	db *gorm.DB
}

// NewCommentService 创建评论服务实例
func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

// CreateComment 创建评论
func (s *CommentService) CreateComment(comment *models.Comment) error {
	return s.db.Create(comment).Error
}

// GetCommentsByPostID 根据文章ID获取评论列表
func (s *CommentService) GetCommentsByPostID(postID uint) ([]models.Comment, error) {
	var comments []models.Comment
	result := s.db.Preload("User").Where("post_id = ?", postID).Order("created_at asc").Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

// GetCommentByID 根据ID获取评论
func (s *CommentService) GetCommentByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	result := s.db.Preload("User").First(&comment, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &comment, nil
}

// GetCommentsByUserID 根据用户ID获取评论列表
func (s *CommentService) GetCommentsByUserID(userID uint) ([]models.Comment, error) {
	var comments []models.Comment
	result := s.db.Preload("User").Preload("Post").Where("user_id = ?", userID).Order("created_at desc").Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

// UpdateComment 更新评论
func (s *CommentService) UpdateComment(comment *models.Comment) error {
	return s.db.Save(comment).Error
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(id uint) error {
	return s.db.Delete(&models.Comment{}, id).Error
}

// GetCommentCountByPostID 获取文章评论数量
func (s *CommentService) GetCommentCountByPostID(postID uint) (int64, error) {
	var count int64
	result := s.db.Model(&models.Comment{}).Where("post_id = ?", postID).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// GetPostByID 根据ID获取文章（用于验证文章是否存在）
func (s *CommentService) GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	result := s.db.First(&post, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}
