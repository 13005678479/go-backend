package models

import (
	"time"
)

// Comment 模型：评论（属于 Post）
type Comment struct {
	ID        uint      `gorm:"primary_key;auto_increment" json:"id"` // 主键自增
	Content   string    `gorm:"type:text;not null" json:"content"`    // 评论内容（非空）
	PostID    uint      `gorm:"not null" json:"post_id"`              // 外键：关联文章 ID
	Post      Post      `gorm:"foreignKey:PostID" json:"post"`        // 所属文章（多对一反向关联）
	UserID    uint      `gorm:"not null" json:"user_id"`             // 外键：关联用户 ID
	User      User      `gorm:"foreignKey:UserID" json:"user"`         // 所属用户（多对一反向关联）
	CreatedAt time.Time `json:"created_at"`                           // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                           // 更新时间
}

// 为模型指定表名
func (Comment) TableName() string { return "biz_comments" }