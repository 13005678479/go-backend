package services

import (
	"blogV2/internal/models"
	"blogV2/pkg/utils"

	"gorm.io/gorm"
)

// PostService 文章服务
type PostService struct {
	db *gorm.DB
}

// NewPostService 创建文章服务实例
func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

// CreatePost 创建文章
func (s *PostService) CreatePost(post *models.Post) error {
	return s.db.Create(post).Error
}

// ListAllPosts 获取所有文章列表
func (s *PostService) ListAllPosts() ([]models.Post, error) {
	var posts []models.Post
	result := s.db.Preload("User").Order("created_at desc").Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

// GetPostByID 根据ID获取文章
func (s *PostService) GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	result := s.db.Preload("User").Preload("Comments").First(&post, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

// GetPostsByUserID 根据用户ID获取文章列表
func (s *PostService) GetPostsByUserID(userID uint) ([]models.Post, error) {
	var posts []models.Post
	result := s.db.Preload("User").Where("user_id = ?", userID).Order("created_at desc").Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

// UpdatePost 更新文章
func (s *PostService) UpdatePost(post *models.Post) error {
	return s.db.Save(post).Error
}

// DeletePost 删除文章
func (s *PostService) DeletePost(id uint) error {
	return s.db.Delete(&models.Post{}, id).Error
}

// GetMostCommentedPost 查询评论数量最多的文章
func (s *PostService) GetMostCommentedPost() (*models.Post, error) {
	var post models.Post
	// 按评论数量降序排序，取第一条
	result := s.db.Preload("User").Preload("Comments").Order("comment_count desc").First(&post)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

// GetPostsWithComments 获取指定用户的所有文章及对应评论
func (s *PostService) GetPostsWithComments(userID uint) ([]models.Post, error) {
	var posts []models.Post
	// 预加载用户信息和评论信息
	result := s.db.Preload("User").Preload("Comments").Where("user_id = ?", userID).Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}