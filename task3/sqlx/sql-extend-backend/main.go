package main

import (
	"log"
	"mysql/database"
	"mysql/models"
	"mysql/services"
)

func main() {
	// 初始化数据库连接
	db := database.InitDB(&models.Account{}, &models.Transaction{})

	// 3. 示例：初始化两个测试账户（可选）
	db.Create(&models.Account{ID: 1, Balance: 1000}) // 账户1：余额1000
	db.Create(&models.Account{ID: 2, Balance: 500})  // 账户2：余额500

	// 4. 执行转账：账户1 向 账户2 转账 300
	err := services.Transfer(db, 1, 2, 300)
	if err != nil {
		log.Fatalf("转账失败：%v", err)
	}
	log.Println("转账成功！")

	// 5. 验证结果：查询两个账户的最新余额
	var acc1, acc2 models.Account
	db.First(&acc1, 1)
	db.First(&acc2, 2)
	log.Printf("账户1余额：%d，账户2余额：%d", acc1.Balance, acc2.Balance)

	// 6. 验证转账记录
	var txn models.Transaction
	db.Order("id desc").First(&txn) // 查询最新一条转账记录
	log.Printf("最新转账记录：转出账户%d → 转入账户%d，金额：%d", txn.FromAccountID, txn.ToAccountID, txn.Amount)
}
