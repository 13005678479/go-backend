// main 包是应用程序的入口点
package main

// 导入所需的包
import (
	"blogV2/config"               // 配置管理包，用于加载应用程序配置
	"blogV2/internal/controllers" // 控制器包，包含业务逻辑处理
	"blogV2/internal/services"    // 服务层包，包含业务逻辑实现
	"blogV2/middleware"           // 中间件包，包含日志和认证中间件
	"blogV2/pkg/database"         // 数据库包，包含数据库连接和操作
	"blogV2/router"               // 路由包，包含路由设置
	"log"                         // Go标准日志包

	"github.com/gin-gonic/gin" // Gin Web框架
)

// main 函数是应用程序的入口点
func main() {
	// 加载应用程序配置
	// 从配置文件或环境变量中读取配置信息
	cfg := config.LoadConfig()

	// 初始化日志中间件
	// 使用配置中的日志设置创建全局日志实例
	middleware.InitLogger(config.GetLogConfig())

	// 记录应用程序启动日志
	middleware.Info("应用程序启动")

	// 设置Gin框架的运行模式
	// 可以是 "debug", "release", "test" 模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库连接
	// 创建与MySQL数据库的连接池
	db := database.InitDB()

	// 自动迁移数据库表结构
	// 根据模型定义自动创建或更新数据库表
	database.AutoMigrate(db)

	// 创建独立的服务层实例
	// 每个服务负责特定的业务领域
	userService := services.NewUserService(db)       // 用户服务，处理用户相关业务
	postService := services.NewPostService(db)       // 文章服务，处理文章相关业务
	commentService := services.NewCommentService(db) // 评论服务，处理评论相关业务

	// 创建控制器层实例
	// 控制器负责接收HTTP请求和返回响应
	userController := controllers.NewUserController(userService)          // 用户控制器
	postController := controllers.NewPostController(postService)          // 文章控制器
	commentController := controllers.NewCommentController(commentService) // 评论控制器

	// 设置HTTP路由
	// 将URL路径映射到对应的控制器方法
	r := router.SetupRouter(userController, postController, commentController)

	// 启动HTTP服务器
	// 记录服务器启动信息
	middleware.Info("服务器启动在端口: %s", cfg.Server.Port)
	middleware.Info("Swagger 文档地址: http://localhost:%s/swagger/index.html", cfg.Server.Port)

	// 启动Gin服务器并监听指定端口
	// 如果启动失败，记录错误并退出程序
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
