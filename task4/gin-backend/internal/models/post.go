package models

import (
	"time"
)

// Post 模型：文章（属于 User，一对多关联 Comment）
type Post struct {
	ID            uint      `gorm:"primary_key;auto_increment" json:"id"`    // 主键自增
	Title         string    `gorm:"type:varchar(200);not null" json:"title"` // 文章标题（非空）
	Content       string    `gorm:"type:text;not null" json:"content"`       // 文章内容（非空）
	UserID        uint      `gorm:"not null" json:"user_id"`                 // 外键：关联用户 ID
	User          User      `gorm:"foreignKey:UserID" json:"user"`           // 所属用户（多对一反向关联）
	CommentCount  int       `gorm:"default:0" json:"comment_count"`          // 评论数量
	CommentStatus string    `gorm:"type:varchar(20);default:'无评论'" json:"comment_status"` // 评论状态
	CreatedAt     time.Time `json:"created_at"`                              // 创建时间
	UpdatedAt     time.Time `json:"updated_at"`                              // 更新时间
	Comments      []Comment `gorm:"foreignKey:PostID" json:"comments"`       // 关联的评论（一对多）
}

// 为模型指定表名
func (Post) TableName() string { return "biz_posts" }