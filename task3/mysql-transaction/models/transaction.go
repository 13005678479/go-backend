package models

// Transaction 对应 transactions 表（转账记录表）
type Transaction struct {
	ID            int64 `gorm:"primary_key;auto_increment"` // 主键自增
	FromAccountID int64 `gorm:"not null"`                   // 转出账户ID（关联 biz_account.id）
	ToAccountID   int64 `gorm:"not null"`                   // 转入账户ID（关联 biz_account.id）
	Amount        int64 `gorm:"type:bigint;not null"`       // 转账金额（正数）
}

// 指定表名
func (Transaction) TableName() string {
	return "biz_transaction"
}
