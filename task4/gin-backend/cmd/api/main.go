package main

import (
	"blogV2/config"
	"blogV2/internal/controllers"
	"blogV2/internal/services"
	"blogV2/pkg/database"
	"blogV2/pkg/utils"
	"blogV2/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.Info("应用程序启动")

	// 加载配置
	cfg := config.LoadConfig()

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库连接
	db := database.InitDB()

	// 自动迁移数据库表
	database.AutoMigrate(db)

	// 创建独立的服务实例
	userService := services.NewUserService(db)
	postService := services.NewPostService(db)
	commentService := services.NewCommentService(db)

	// 创建控制器实例
	userController := controllers.NewUserController(userService)
	postController := controllers.NewPostController(postService)
	commentController := controllers.NewCommentController(commentService)

	// 设置路由
	r := router.SetupRouter(userController, postController, commentController)

	// 启动服务器
	utils.Info("服务器启动在端口: %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}