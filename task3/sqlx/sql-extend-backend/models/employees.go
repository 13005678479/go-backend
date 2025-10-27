package models

// Account 对应 accounts 表（账户表）
type Employee struct {
	ID      int64 `gorm:"primary_key;auto_increment"` // 主键自增
	Balance int64 `gorm:"type:bigint;not null"`       // 账户余额（bigint 避免溢出）
}

// 指定表名
func (Employee) TableName() string {
	return "biz_employee"
}
