package models

// 包含字段 ID 、 Name 、 Department 、 Salary
type Employee struct {
	ID         int64  `gorm:"primary_key;auto_increment"` // 主键自增
	Name       string `gorm:"type:varchar(255)"`          // 姓名（字符串类型）
	Department string `gorm:"type:varchar(255)"`          // 部门
	Salary     int32  `gorm:"type:bigint;not null"`       // salary
}

// 指定表名
func (Employee) TableName() string {
	return "biz_employee"
}
