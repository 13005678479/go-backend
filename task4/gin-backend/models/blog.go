package main

import (
	"time"
)

// User 模型：用户（一对多关联 Post）
type User struct {
	ID        uint      `gorm:"primary_key;auto_increment" json:"id"`             // 主键自增
	Username  string    `gorm:"type:varchar(50);unique;not null" json:"username"` // 用户名（唯一、非空）
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email"`   // 邮箱（唯一、非空）
	Password  string    `gorm:"type:varchar(100);not null" json:"-"`              // 密码（不返回给前端）
	CreatedAt time.Time `json:"created_at"`                                       // 创建时间（GORM 自动维护）
	UpdatedAt time.Time `json:"updated_at"`                                       // 更新时间（GORM 自动维护）
	Posts     []Post    `gorm:"foreignKey:UserID" json:"posts"`                   // 关联的文章（一对多）
}

// Post 模型：文章（属于 User，一对多关联 Comment）
type Post struct {
	ID        uint      `gorm:"primary_key;auto_increment" json:"id"`    // 主键自增
	Title     string    `gorm:"type:varchar(200);not null" json:"title"` // 文章标题（非空）
	Content   string    `gorm:"type:text;not null" json:"content"`       // 文章内容（非空）
	UserID    uint      `gorm:"not null" json:"user_id"`                 // 外键：关联用户 ID
	User      User      `gorm:"foreignKey:UserID" json:"user"`           // 所属用户（多对一反向关联）
	CreatedAt time.Time `json:"created_at"`                              // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                              // 更新时间
	Comments  []Comment `gorm:"foreignKey:PostID" json:"comments"`       // 关联的评论（一对多）
}

// Comment 模型：评论（属于 Post）
type Comment struct {
	ID        uint      `gorm:"primary_key;auto_increment" json:"id"` // 主键自增
	Content   string    `gorm:"type:text;not null" json:"content"`    // 评论内容（非空）
	PostID    uint      `gorm:"not null" json:"post_id"`              // 外键：关联文章 ID
	Post      Post      `gorm:"foreignKey:PostID" json:"post"`        // 所属文章（多对一反向关联）
	CreatedAt time.Time `json:"created_at"`                           // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                           // 更新时间
}

// 为模型指定表名（可选，GORM 默认使用复数形式，这里显式指定更清晰）
func (User) TableName() string    { return "biz_users" }
func (Post) TableName() string    { return "biz_posts" }
func (Comment) TableName() string { return "biz_comments" }
