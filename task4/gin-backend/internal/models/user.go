package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 模型：用户（一对多关联 Post）
type User struct {
	ID        uint      `gorm:"primary_key;auto_increment" json:"id"`             // 主键自增
	Username  string    `gorm:"type:varchar(50);unique;not null" json:"username"` // 用户名（唯一、非空）
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email"`   // 邮箱（唯一、非空）
	Password  string    `gorm:"type:varchar(100);not null" json:"-"`              // 密码（不返回给前端）
	PostCount int       `gorm:"default:0" json:"post_count"`                     // 文章数量
	CreatedAt time.Time `json:"created_at"`                                       // 创建时间（GORM 自动维护）
	UpdatedAt time.Time `json:"updated_at"`                                       // 更新时间（GORM 自动维护）
	Posts     []Post    `gorm:"foreignKey:UserID" json:"posts"`                   // 关联的文章（一对多）
}

// 为模型指定表名
func (User) TableName() string { return "biz_users" }

// BeforeSave User保存前钩子：密码加密
func (u *User) BeforeSave(tx *gorm.DB) error {
	if len(u.Password) > 0 && u.Password[0] != '$' { // 检查是否已经加密
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}