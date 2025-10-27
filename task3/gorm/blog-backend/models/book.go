package models

// books 表，包含字段 id 、 title 、 author 、 price
type Book struct {
	ID     int64   `gorm:"primary_key;auto_increment" json:"id"`     // 主键自增，映射 id 字段
	Title  string  `gorm:"type:varchar(255);not null" json:"title"`  // 书名，非空字符串，映射 title 字段
	Author string  `gorm:"type:varchar(255);not null" json:"author"` // 作者，非空字符串，映射 author 字段
	Price  float64 `gorm:"type:decimal(10,2);not null" json:"price"` // 价格，精确到2位小数，映射 price 字段
}

// 指定表名
func (Book) TableName() string {
	return "biz_book"
}
