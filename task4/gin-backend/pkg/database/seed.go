package database

import (
	"blogV2/internal/models"
	"blogV2/middleware"
	"time"

	"gorm.io/gorm"
)

// SeedData 初始化测试数据
func SeedData(db *gorm.DB) {
	middleware.Info("开始初始化测试数据...")

	// 检查是否已有数据
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	
	if userCount > 0 {
		middleware.Info("数据库中已有数据，跳过初始化")
		return
	}

	// 创建测试用户
	users := []models.User{
		{
			Username: "admin",
			Email:    "admin@example.com",
			Password: "admin123",
		},
		{
			Username: "alice",
			Email:    "alice@example.com",
			Password: "alice123",
		},
		{
			Username: "bob",
			Email:    "bob@example.com",
			Password: "bob123",
		},
		{
			Username: "charlie",
			Email:    "charlie@example.com",
			Password: "charlie123",
		},
	}

	// 批量创建用户
	for i := range users {
		result := db.Create(&users[i])
		if result.Error != nil {
			middleware.Error("创建用户失败: %v", result.Error)
			continue
		}
		middleware.Info("创建用户: %s (ID: %d)", users[i].Username, users[i].ID)
	}

	// 创建测试文章
	posts := []models.Post{
		{
			Title:         "欢迎来到博客系统",
			Content:       "这是第一篇测试文章，欢迎使用我们的博客系统！",
			UserID:        users[0].ID,
			CommentStatus: "开放评论",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			Title:         "Go语言开发指南",
			Content:       "Go语言是一种高效、简洁的编程语言，特别适合后端开发。",
			UserID:        users[1].ID,
			CommentStatus: "开放评论",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			Title:         "Gin框架入门教程",
			Content:       "Gin是Go语言中最流行的Web框架之一，具有高性能和易用性。",
			UserID:        users[2].ID,
			CommentStatus: "开放评论",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			Title:         "数据库设计最佳实践",
			Content:       "良好的数据库设计是系统稳定性的基础，本文分享一些设计经验。",
			UserID:        users[3].ID,
			CommentStatus: "开放评论",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			Title:         "RESTful API设计原则",
			Content:       "遵循RESTful原则设计API可以提高系统的可维护性和扩展性。",
			UserID:        users[0].ID,
			CommentStatus: "开放评论",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	// 批量创建文章
	for i := range posts {
		result := db.Create(&posts[i])
		if result.Error != nil {
			middleware.Error("创建文章失败: %v", result.Error)
			continue
		}
		middleware.Info("创建文章: %s (ID: %d)", posts[i].Title, posts[i].ID)
	}

	// 创建测试评论
	comments := []models.Comment{
		{
			Content:   "这篇文章写得很好，对我帮助很大！",
			PostID:    posts[0].ID,
			UserID:    users[1].ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Content:   "期待更多关于Go语言的教程！",
			PostID:    posts[1].ID,
			UserID:    users[2].ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Content:   "Gin框架确实很强大，感谢分享！",
			PostID:    posts[2].ID,
			UserID:    users[3].ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Content:   "数据库设计经验很实用，学习了！",
			PostID:    posts[3].ID,
			UserID:    users[0].ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Content:   "RESTful API设计确实很重要，谢谢分享！",
			PostID:    posts[4].ID,
			UserID:    users[1].ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Content:   "希望能看到更多实战案例！",
			PostID:    posts[0].ID,
			UserID:    users[2].ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// 批量创建评论
	for i := range comments {
		result := db.Create(&comments[i])
		if result.Error != nil {
			middleware.Error("创建评论失败: %v", result.Error)
			continue
		}
		middleware.Info("创建评论: %s (ID: %d)", comments[i].Content[:20]+"...", comments[i].ID)
	}

	// 更新用户和文章的计数
	for i := range users {
		var postCount int64
		db.Model(&models.Post{}).Where("user_id = ?", users[i].ID).Count(&postCount)
		db.Model(&models.User{}).Where("id = ?", users[i].ID).Update("post_count", postCount)
	}

	for i := range posts {
		var commentCount int64
		db.Model(&models.Comment{}).Where("post_id = ?", posts[i].ID).Count(&commentCount)
		db.Model(&models.Post{}).Where("id = ?", posts[i].ID).Update("comment_count", commentCount)
	}

	middleware.Info("测试数据初始化完成")
	middleware.Info("创建了 %d 个用户，%d 篇文章，%d 条评论", len(users), len(posts), len(comments))
}