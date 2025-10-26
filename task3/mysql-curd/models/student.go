package models

type Student struct {
	ID    int    `gorm:"primary_key;auto_increment"` // 主键且自增
	Name  string `gorm:"type:varchar(255)"`          // 学生姓名（字符串类型）
	Age   int    // 学生年龄（整数类型）
	Grade string `gorm:"type:varchar(50)"` // 学生年级（字符串类型）
}

// 指定表名
func (Student) TableName() string {
	return "biz_students"
}
