package api

import (
	"blog/pkg/database"
	"blog/pkg/utils"
)

func main() {
	utils.Info("应用程序启动")

	// 初始化数据库连接
	database.InitDB()

	// // 示例：插入测试数据（可选，用于验证查询）
	// testBooks := []models.Book{
	// 	{Title: "Go 编程实战", Author: "张三", Price: 69.90},
	// 	{Title: "Python 入门", Author: "李四", Price: 45.50},
	// 	{Title: "Java 高级开发", Author: "王五", Price: 89.00},
	// }

	// utils.Info("开始插入测试数据")
	// result := db.Create(&testBooks) // 批量插入测试数据
	// if result.Error != nil {
	// 	utils.Error("插入测试数据失败: %v", result.Error)
	// } else {
	// 	utils.Info("成功插入 %d 条测试数据", len(testBooks))
	// }

	// svc := services.New(db)

	// // 执行查询：价格大于 50 元的书籍
	// expensiveBooks, err := svc.GetBooksPriceGreaterThan50()
	// if err != nil {
	// 	utils.Error("查询失败: %v", err)
	// } else {
	// 	utils.Info("查询结果：价格大于50元的书籍列表")
	// 	fmt.Println("价格大于 50列表：")
	// 	for _, book := range expensiveBooks {
	// 		fmt.Printf("ID: %d, 书名: %s, 作者: %s, 价格: %.2f\n", book.ID, book.Title, book.Author, book.Price)
	// 	}
	// }

	utils.Info("应用程序正常结束")
}
