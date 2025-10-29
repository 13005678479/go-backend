// router 包负责设置应用程序的HTTP路由
package router

// 导入所需的包
import (
	"blogV2/docs"                    // Swagger文档包，包含API文档信息
	"blogV2/internal/controllers"     // 控制器包，包含业务逻辑处理
	"blogV2/middleware"              // 中间件包，包含认证等中间件

	"github.com/gin-gonic/gin"        // Gin Web框架
	swaggerFiles "github.com/swaggo/files"     // Swagger静态文件处理器
	ginSwagger "github.com/swaggo/gin-swagger" // Gin Swagger中间件
)

// SetupRouter 函数设置应用程序的所有HTTP路由
// 参数:
//   userController - 用户控制器实例
//   postController - 文章控制器实例  
//   commentController - 评论控制器实例
// 返回值:
//   *gin.Engine - 配置好的Gin引擎实例
func SetupRouter(
	userController *controllers.UserController,
	postController *controllers.PostController,
	commentController *controllers.CommentController,
) *gin.Engine {
	// 创建默认的Gin引擎实例
	// 包含日志、恢复等默认中间件
	r := gin.Default()

	// 应用全局日志中间件
	// 记录所有HTTP请求的详细信息
	r.Use(middleware.LoggerMiddleware(middleware.GetLogger()))

	// 配置Swagger文档路由
	// 设置API的基础路径为"/api"
	docs.SwaggerInfo.BasePath = "/api"
	// 设置Swagger UI路由，可以访问/swagger/index.html查看API文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 设置公开路由组 - 不需要认证即可访问的路由
	// 所有以"/api"开头的路径都属于这个路由组
	public := r.Group("/api")
	{
		// 用户相关公开路由
		public.POST("/register", userController.Register) // 用户注册接口
		public.POST("/login", userController.Login)       // 用户登录接口

		// 文章相关公开路由 - 所有人都可以访问
		public.GET("/posts", postController.ListPosts)     // 获取文章列表
		public.GET("/posts/:id", postController.GetPost)   // 获取指定文章详情
	}

	// 设置需要认证的路由组 - 需要JWT token认证
	// 所有以"/api"开头的路径都需要经过认证中间件
	protected := r.Group("/api")
	// 应用认证中间件，验证JWT token
	protected.Use(middleware.AuthMiddleware())
	{
		// 用户相关受保护路由
		protected.GET("/user/me", userController.GetCurrentUser) // 获取当前登录用户信息

		// 文章相关受保护路由 - 需要登录才能操作
		protected.POST("/posts", postController.CreatePost)      // 创建新文章
		protected.PUT("/posts/:id", postController.UpdatePost)   // 更新指定文章
		protected.DELETE("/posts/:id", postController.DeletePost) // 删除指定文章

		// 评论相关受保护路由
		protected.POST("/posts/:id/comments", commentController.CreateComment) // 为文章创建评论
		protected.GET("/posts/:id/comments", commentController.ListComments)  // 获取文章的所有评论
		protected.DELETE("/comments/:id", commentController.DeleteComment)   // 删除指定评论
	}

	// 返回配置好的路由引擎
	return r
}
