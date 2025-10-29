package router

import (
	"blogV2/docs"
	"blogV2/internal/controllers"
	"blogV2/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(
	userController *controllers.UserController,
	postController *controllers.PostController,
	commentController *controllers.CommentController,
) *gin.Engine {
	r := gin.Default()

	// Swagger 文档路由
	docs.SwaggerInfo.BasePath = "/api"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 公开路由
	public := r.Group("/api")
	{
		// 用户相关
		public.POST("/register", userController.Register)
		public.POST("/login", userController.Login)

		// 文章相关（公开访问）
		public.GET("/posts", postController.ListPosts)
		public.GET("/posts/:id", postController.GetPost)
	}

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// 用户相关
		protected.GET("/user/me", userController.GetCurrentUser)

		// 文章相关（需要认证）
		protected.POST("/posts", postController.CreatePost)
		protected.PUT("/posts/:id", postController.UpdatePost)
		protected.DELETE("/posts/:id", postController.DeletePost)

		// 评论相关
		protected.POST("/posts/:id/comments", commentController.CreateComment)
		protected.GET("/posts/:id/comments", commentController.ListComments)
		protected.DELETE("/comments/:id", commentController.DeleteComment)
	}

	return r
}
