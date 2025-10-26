package services

import "gorm.io/gorm"

// Transfer 转账函数：从 fromID 账户向 toID 账户转账 amount 金额
func Transfer(db *gorm.DB, fromID, toID, amount int64) {
	// db.Transaction()

}
