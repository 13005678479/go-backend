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
func getUserPostsWithComments(db *gorm.DB, userID uint) ([]models.Post, error) {
	var posts []models.Post
	// 预加载用户信息和评论信息
	result := db.Preload("User").Preload("Comments").Where("user_id = ?", userID).Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

// GetMostCommentedPost 查询评论数量最多的文章
func GetMostCommentedPost(db *gorm.DB) (models.Post, error) {
	var post models.Post
	// 按评论数量降序排序，取第一条
	result := db.Preload("User").Preload("Comments").Order("comments_count desc").First(&post)
	if result.Error != nil {
		return models.Post{}, result.Error
	}
	return post, nil
}

// 3. 钩子函数逻辑

// BeforeCreate Post创建前钩子：更新用户文章数量
func (p *models.Post) BeforeCreate(tx *gorm.DB) error {
	// 自增用户的文章数量
	return tx.Model(&models.User{}).Where("id = ?", p.UserID).Update("post_count", gorm.Expr("post_count + ?", 1)).Error
}

// AfterCreate Comment创建后钩子：更新文章评论数量和状态
func (c *models.Comment) AfterCreate(tx *gorm.DB) error {
	// 自增文章的评论数量
	if err := tx.Model(&Post{}).Where("id = ?", c.PostID).
		Updates(map[string]interface{}{
			"comment_count":  gorm.Expr("comment_count + ?", 1),
			"comment_status": "有评论",
		}).Error; err != nil {
		return err
	}
	return nil
}

// AfterDelete Comment删除后钩子：检查文章评论状态
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// 先查询当前文章的评论数量
	var post Post
	if err := tx.First(&post, c.PostID).Error; err != nil {
		return err
	}

	// 自减评论数量
	newCount := post.CommentCount - 1
	updateData := map[string]interface{}{"comment_count": newCount}

	// 如果评论数量为0，更新状态为"无评论"
	if newCount == 0 {
		updateData["comment_status"] = "无评论"
	}

	return tx.Model(&Post{}).Where("id = ?", c.PostID).Updates(updateData).Error
}
