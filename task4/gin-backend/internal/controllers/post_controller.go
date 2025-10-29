// controllers 包包含所有HTTP请求处理器的实现
// 这个文件专门处理文章相关的HTTP请求
package controllers

// 导入所需的包
import (
	"blogV2/internal/models"     // 数据模型包，包含数据库实体定义
	"blogV2/internal/services"   // 服务层包，包含业务逻辑实现
	"blogV2/middleware"          // 中间件包，包含日志记录等中间件
	"blogV2/pkg/utils"          // 工具包，包含工具函数
	"net/http"                   // HTTP协议包，包含HTTP状态码等常量
	"strconv"                    // 字符串转换包，用于字符串和数字的转换

	"github.com/gin-gonic/gin"    // Gin Web框架
	"gorm.io/gorm"                // GORM ORM框架，用于数据库操作
)

// PostController 文章控制器
// @Summary 文章管理接口
// @Description 提供文章的创建、查询、更新、删除等操作
// @Tags posts
// @Accept json
// @Produce json
// @Router /api/v1/posts [get]

// PostController 文章控制器结构体
// 负责处理所有与文章相关的HTTP请求
// 遵循MVC架构模式，作为视图层和模型层之间的桥梁
type PostController struct {
	service *services.PostService // 文章服务实例，用于处理业务逻辑
	                           // 通过依赖注入的方式注入，实现控制层与服务层的分离
}

// NewPostController 创建文章控制器实例
// 这是控制器的工厂函数，用于创建并初始化PostController实例
// 参数:
//   service - 文章服务实例，包含文章相关的业务逻辑
// 返回值:
//   *PostController - 初始化好的文章控制器实例
func NewPostController(service *services.PostService) *PostController {
	// 返回新创建的PostController实例
	// 将传入的服务实例赋值给控制器的service字段
	return &PostController{service: service}
}

// CreatePost 创建文章
// @Summary 创建文章
// @Description 创建新的文章
// @Tags posts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreatePostRequest true "文章信息"
// @Success 201 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/posts [post]

// CreatePostRequest 创建文章请求结构体
// 定义了客户端发送创建文章请求时需要提供的字段和验证规则
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=200"` // 文章标题，必填字段，长度限制3-200个字符
	Content string `json:"content" binding:"required,min=10"`      // 文章内容，必填字段，最少10个字符
}

// CreatePost 处理创建文章请求
// 这是HTTP POST /api/posts 路由的处理函数（需要认证）
// 参数:
//   ctx - Gin上下文对象，包含HTTP请求和响应的所有信息
func (c *PostController) CreatePost(ctx *gin.Context) {
	// 从Gin上下文中获取当前登录用户的ID
	// 这个值由认证中间件在验证JWT token后设置
	userID, _ := ctx.Get("userID")

	// 声明创建文章请求变量，用于存储解析后的请求数据
	var req CreatePostRequest
	
	// 使用Gin框架的ShouldBindJSON方法将请求的JSON体绑定到CreatePostRequest结构体
	// 这个方法会自动验证binding标签定义的规则
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 如果绑定或验证失败，记录错误日志并返回400 Bad Request状态码
		middleware.Error("创建文章参数验证失败: " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return // 提前返回，不再执行后续代码
	}

	// 创建文章模型实例，将请求数据转换为数据库实体
	// 这个实例将被传递给服务层进行业务处理
	post := &models.Post{
		Title:   req.Title,     // 设置文章标题，从请求中获取
		Content: req.Content,  // 设置文章内容，从请求中获取
		UserID:  userID.(uint), // 设置作者ID，从认证中间件中获取
	}

	// 调用文章服务的CreatePost方法创建新文章
	// 这个方法会处理数据验证、数据库插入等业务逻辑
	if err := c.service.CreatePost(post); err != nil {
		// 如果创建文章过程中发生错误，记录错误日志并返回500 Internal Server Error状态码
		middleware.Error("创建文章失败 (用户ID: " + utils.UintToString(userID.(uint)) + "): " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return // 提前返回，不再执行后续代码
	}

	// 文章创建成功，记录信息日志并返回201 Created状态码
	// 按照RESTful API最佳实践，创建资源成功应该返回201状态码
	middleware.Info("文章创建成功 (ID: " + strconv.Itoa(int(post.ID)) + ", 用户: " + utils.UintToString(userID.(uint)) + ")")
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "文章创建成功", // 成功消息，告知客户端创建成功
		"post": gin.H{ // 返回创建的文章信息
			"id":         post.ID,       // 文章ID，由数据库自动生成
			"title":      post.Title,   // 文章标题
			"content":    post.Content, // 文章内容
			"created_at": post.CreatedAt.Format("2006-01-02 15:04:05"), // 格式化后的创建时间
		},
	})
}

// GetPost 获取文章详情
// @Summary 获取文章详情
// @Description 根据文章ID获取文章详细信息
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} map[string]interface{} "文章详情"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "文章不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/posts/{id} [get]

// GetPost 处理获取文章详情请求
// 这是HTTP GET /api/posts/:id 路由的处理函数（公开访问）
// 参数:
//   ctx - Gin上下文对象，包含HTTP请求和响应的所有信息
func (c *PostController) GetPost(ctx *gin.Context) {
	// 从URL路径参数中获取文章ID字符串
	postIDStr := ctx.Param("id")
	
	// 将字符串类型的文章ID转换为整数类型
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		// 如果转换失败，说明文章ID格式不正确，记录错误日志并返回400 Bad Request状态码
		middleware.Error("文章ID格式错误: " + postIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return // 提前返回，不再执行后续代码
	}

	// 调用文章服务的GetPostByID方法获取文章详情
	// 这个方法会查询数据库并返回文章信息（包含作者信息）
	post, err := c.service.GetPostByID(uint(postID))
	if err != nil {
		// 如果查询过程中发生错误
		if err == gorm.ErrRecordNotFound {
			// 如果错误是记录不存在，返回404 Not Found状态码
			ctx.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
			return
		}
		// 其他错误返回500 Internal Server Error状态码
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return // 提前返回，不再执行后续代码
	}

	// 文章获取成功，返回200 OK状态码和文章详情
	ctx.JSON(http.StatusOK, gin.H{
		"post": gin.H{ // 文章详细信息
			"id":         post.ID,       // 文章ID
			"title":      post.Title,   // 文章标题
			"content":    post.Content, // 文章内容
			"user_id":    post.UserID,  // 作者ID
			"username":   post.User.Username, // 作者用户名
			"created_at": post.CreatedAt.Format("2006-01-02 15:04:05"), // 格式化后的创建时间
			"updated_at": post.UpdatedAt.Format("2006-01-02 15:04:05"), // 格式化后的更新时间
		},
	})
}

// UpdatePost 更新文章
// @Summary 更新文章
// @Description 更新指定文章的内容
// @Tags posts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "文章ID"
// @Param request body CreatePostRequest true "更新后的文章信息"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 403 {object} map[string]interface{} "无权操作"
// @Failure 404 {object} map[string]interface{} "文章不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/posts/{id} [put]

// UpdatePost 处理更新文章请求
// 这是HTTP PUT /api/posts/:id 路由的处理函数（需要认证）
// 参数:
//   ctx - Gin上下文对象，包含HTTP请求和响应的所有信息
func (c *PostController) UpdatePost(ctx *gin.Context) {
	// 从Gin上下文中获取当前登录用户的ID
	userID, _ := ctx.Get("userID")
	
	// 从URL路径参数中获取文章ID字符串
	postIDStr := ctx.Param("id")
	
	// 将字符串类型的文章ID转换为整数类型
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		// 如果转换失败，说明文章ID格式不正确，记录错误日志并返回400 Bad Request状态码
		middleware.Error("文章ID格式错误: " + postIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return // 提前返回，不再执行后续代码
	}

	// 声明更新文章请求变量，用于存储解析后的请求数据
	var req CreatePostRequest
	
	// 使用Gin框架的ShouldBindJSON方法将请求的JSON体绑定到CreatePostRequest结构体
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 如果绑定或验证失败，记录错误日志并返回400 Bad Request状态码
		middleware.Error("更新文章参数验证失败: " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return // 提前返回，不再执行后续代码
	}

	// 验证文章是否存在且属于当前用户
	// 先获取文章详情，检查权限
	post, err := c.service.GetPostByID(uint(postID))
	if err != nil {
		// 如果查询过程中发生错误
		if err == gorm.ErrRecordNotFound {
			// 如果错误是记录不存在，返回404 Not Found状态码
			ctx.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
			return
		}
		// 其他错误返回500 Internal Server Error状态码
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return // 提前返回，不再执行后续代码
	}

	// 检查当前用户是否有权限更新这篇文章
	// 只有文章的作者才能更新自己的文章
	if post.UserID != userID.(uint) {
		// 如果当前用户不是文章作者，记录警告日志并返回403 Forbidden状态码
		middleware.Warn("无权更新文章 (用户ID: " + utils.UintToString(userID.(uint)) + ", 文章ID: " + postIDStr + ")")
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权更新此文章"})
		return // 提前返回，不再执行后续代码
	}

	// 更新文章内容
	// 使用请求中的新数据更新文章对象
	post.Title = req.Title   // 更新文章标题
	post.Content = req.Content // 更新文章内容
	
	// 调用文章服务的UpdatePost方法更新文章
	// 这个方法会更新数据库中的文章记录
	if err := c.service.UpdatePost(post); err != nil {
		// 如果更新过程中发生错误，记录错误日志并返回500 Internal Server Error状态码
		middleware.Error("更新文章失败 (ID: " + postIDStr + "): " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return // 提前返回，不再执行后续代码
	}

	// 文章更新成功，记录信息日志并返回200 OK状态码
	middleware.Info("文章更新成功 (ID: " + postIDStr + ")")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "文章更新成功", // 成功消息
		"post": gin.H{ // 返回更新后的文章信息
			"id":         post.ID,       // 文章ID
			"title":      post.Title,   // 更新后的文章标题
			"content":    post.Content, // 更新后的文章内容
			"updated_at": post.UpdatedAt.Format("2006-01-02 15:04:05"), // 格式化后的更新时间
		},
	})
}

// DeletePost 删除文章
// @Summary 删除文章
// @Description 删除指定文章
// @Tags posts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "文章ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 403 {object} map[string]interface{} "无权操作"
// @Failure 404 {object} map[string]interface{} "文章不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/posts/{id} [delete]

// DeletePost 处理删除文章请求
// 这是HTTP DELETE /api/posts/:id 路由的处理函数（需要认证）
// 参数:
//   ctx - Gin上下文对象，包含HTTP请求和响应的所有信息
func (c *PostController) DeletePost(ctx *gin.Context) {
	// 从Gin上下文中获取当前登录用户的ID
	userID, _ := ctx.Get("userID")
	
	// 从URL路径参数中获取文章ID字符串
	postIDStr := ctx.Param("id")
	
	// 将字符串类型的文章ID转换为整数类型
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		// 如果转换失败，说明文章ID格式不正确，记录错误日志并返回400 Bad Request状态码
		middleware.Error("文章ID格式错误: " + postIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return // 提前返回，不再执行后续代码
	}

	// 验证文章是否存在且属于当前用户
	// 先获取文章详情，检查权限
	post, err := c.service.GetPostByID(uint(postID))
	if err != nil {
		// 如果查询过程中发生错误
		if err == gorm.ErrRecordNotFound {
			// 如果错误是记录不存在，返回404 Not Found状态码
			ctx.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
			return
		}
		// 其他错误返回500 Internal Server Error状态码
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return // 提前返回，不再执行后续代码
	}

	// 检查当前用户是否有权限删除这篇文章
	// 只有文章的作者才能删除自己的文章
	if post.UserID != userID.(uint) {
		// 如果当前用户不是文章作者，记录警告日志并返回403 Forbidden状态码
		middleware.Warn("无权删除文章 (用户ID: " + utils.UintToString(userID.(uint)) + ", 文章ID: " + postIDStr + ")")
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权删除此文章"})
		return // 提前返回，不再执行后续代码
	}

	// 调用文章服务的DeletePost方法删除文章
	// 这个方法会从数据库中删除文章记录
	if err := c.service.DeletePost(uint(postID)); err != nil {
		// 如果删除过程中发生错误，记录错误日志并返回500 Internal Server Error状态码
		middleware.Error("删除文章失败 (ID: " + postIDStr + "): " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败"})
		return // 提前返回，不再执行后续代码
	}

	// 文章删除成功，记录信息日志并返回200 OK状态码
	middleware.Info("文章删除成功 (ID: " + postIDStr + ")")
	ctx.JSON(http.StatusOK, gin.H{"message": "文章删除成功"}) // 返回成功消息
}

// ListPosts 获取所有文章列表
// @Summary 获取文章列表
// @Description 获取所有文章的列表
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "文章列表"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/posts [get]

// ListPosts 处理获取文章列表请求
// 这是HTTP GET /api/posts 路由的处理函数（公开访问）
// 参数:
//   ctx - Gin上下文对象，包含HTTP请求和响应的所有信息
func (c *PostController) ListPosts(ctx *gin.Context) {
	// 调用文章服务的ListAllPosts方法获取所有文章列表
	// 这个方法会查询数据库并返回所有文章（包含作者信息）
	posts, err := c.service.ListAllPosts()
	if err != nil {
		// 如果查询过程中发生错误，记录错误日志并返回500 Internal Server Error状态码
		middleware.Error("获取文章列表失败: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return // 提前返回，不再执行后续代码
	}

	// 格式化返回数据
	// 将文章列表转换为前端友好的格式
	var result []gin.H
	for _, post := range posts {
		result = append(result, gin.H{
			"id":         post.ID,       // 文章ID
			"title":      post.Title,   // 文章标题
			"content":    post.Content, // 文章内容
			"user_id":    post.UserID,  // 作者ID
			"username":   post.User.Username, // 作者用户名
			"created_at": post.CreatedAt.Format("2006-01-02 15:04:05"), // 格式化后的创建时间
			"updated_at": post.UpdatedAt.Format("2006-01-02 15:04:05"), // 格式化后的更新时间
		})
	}

	// 返回文章列表，状态码200 OK
	ctx.JSON(http.StatusOK, gin.H{
		"total": len(result), // 文章总数
		"posts": result,     // 文章列表
	})
}