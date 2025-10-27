package services

import (
	"fmt"
	"mysql/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Transfer(db *gorm.DB, fromID, toID, amount int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 查询转出账户（加行锁）
		var fromAcc models.Account
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).First(&fromAcc, fromID).Error; err != nil {
			return err
		}
		// 查询转入账户（加行锁）
		var toAcc models.Account
		if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).First(&toAcc, toID).Error; err != nil {
			return err
		}
		// 校验余额
		if fromAcc.Balance < amount {
			return fmt.Errorf("账户 %d 余额不足，当前余额：%d，转账金额：%d", fromID, fromAcc.Balance, amount)
		}
		// 更新余额
		if err := tx.Model(&fromAcc).Update("balance", fromAcc.Balance-amount).Error; err != nil {
			return err
		}
		if err := tx.Model(&toAcc).Update("balance", toAcc.Balance+amount).Error; err != nil {
			return err
		}
		// 记录转账
		transaction := models.Transaction{
			FromAccountID: fromID,
			ToAccountID:   toID,
			Amount:        amount,
		}
		return tx.Create(&transaction).Error
	})
}
