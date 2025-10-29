// controllers 包包含所有HTTP请求处理器的实现
// 这个文件专门处理评论相关的HTTP请求
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

// CommentController 评论控制器
// @Summary 评论管理接口
// @Description 提供评论的创建、查询、删除等操作
// @Tags comments
// @Accept json
// @Produce json
// @Router /api/v1/comments [get]

// CommentController 评论控制器结构体
// 负责处理所有与评论相关的HTTP请求
// 遵循MVC架构模式，作为视图层和模型层之间的桥梁
type CommentController struct {
	service *services.CommentService // 评论服务实例，用于处理业务逻辑
	                              // 通过依赖注入的方式注入，实现控制层与服务层的分离
}

// NewCommentController 创建评论控制器实例
// 这是控制器的工厂函数，用于创建并初始化CommentController实例
// 参数:
//   service - 评论服务实例，包含评论相关的业务逻辑
// 返回值:
//   *CommentController - 初始化好的评论控制器实例
func NewCommentController(service *services.CommentService) *CommentController {
	// 返回新创建的CommentController实例
	// 将传入的服务实例赋值给控制器的service字段
	return &CommentController{service: service}
}

// CreateComment 创建评论
// @Summary 创建评论
// @Description 为指定文章创建评论
// @Tags comments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "文章ID"
// @Param request body CreateCommentRequest true "评论内容"
// @Success 201 {object} map[string]interface{} "评论创建成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "文章不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/posts/{id}/comments [post]

// CreateCommentRequest 创建评论请求结构体
// 定义了客户端发送创建评论请求时需要提供的字段和验证规则
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=500"` // 评论内容，必填字段，长度限制1-500个字符
}

// CreateComment 处理创建评论请求
// 这是HTTP POST /api/posts/:id/comments 路由的处理函数（需要认证）
// 参数:
//   ctx - Gin上下文对象，包含HTTP请求和响应的所有信息
func (c *CommentController) CreateComment(ctx *gin.Context) {
	// 从Gin上下文中获取当前登录用户的ID
	// 这个值由认证中间件在验证JWT token后设置
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

	// 验证文章是否存在
	// 调用评论服务的GetPostByID方法验证文章是否存在
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

	// 声明创建评论请求变量，用于存储解析后的请求数据
	var req CreateCommentRequest
	
	// 使用Gin框架的ShouldBindJSON方法将请求的JSON体绑定到CreateCommentRequest结构体
	// 这个方法会自动验证binding标签定义的规则
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 如果绑定或验证失败，记录错误日志并返回400 Bad Request状态码
		middleware.Error("创建评论参数验证失败: " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return // 提前返回，不再执行后续代码
	}

	// 创建评论模型实例，将请求数据转换为数据库实体
	// 这个实例将被传递给服务层进行业务处理
	comment := &models.Comment{
		Content: req.Content,    // 设置评论内容，从请求中获取
		PostID:  post.ID,        // 设置文章ID，从验证的文章中获取
		UserID:  userID.(uint),  // 设置评论者ID，从认证中间件中获取
	}

	// 调用评论服务的CreateComment方法创建新评论
	// 这个方法会处理数据验证、数据库插入等业务逻辑
	if err := c.service.CreateComment(comment); err != nil {
		// 如果创建评论过程中发生错误，记录错误日志并返回500 Internal Server Error状态码
		middleware.Error("创建评论失败 (文章ID: " + postIDStr + "): " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败"})
		return // 提前返回，不再执行后续代码
	}

	// 评论创建成功，记录信息日志并返回201 Created状态码
	// 按照RESTful API最佳实践，创建资源成功应该返回201状态码
	middleware.Info("评论创建成功 (ID: " + strconv.Itoa(int(comment.ID)) + ", 文章ID: " + postIDStr + ")")
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "评论创建成功", // 成功消息，告知客户端创建成功
		"comment": gin.H{ // 返回创建的评论信息
			"id":         comment.ID,       // 评论ID，由数据库自动生成
			"content":    comment.Content, // 评论内容
			"created_at": comment.CreatedAt.Format("2006-01-02 15:04:05"), // 格式化后的创建时间
		},
	})
}

// ListComments 获取文章的所有评论
// @Summary 获取文章评论列表
// @Description 获取指定文章的所有评论
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} map[string]interface{} "评论列表"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "文章不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/posts/{id}/comments [get]

// ListComments 处理获取文章评论列表请求
// 这是HTTP GET /api/posts/:id/comments 路由的处理函数（公开访问）
// 参数:
//   ctx - Gin上下文对象，包含HTTP请求和响应的所有信息
func (c *CommentController) ListComments(ctx *gin.Context) {
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

	// 验证文章是否存在
	// 调用评论服务的GetPostByID方法验证文章是否存在
	_, err = c.service.GetPostByID(uint(postID))
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

	// 调用评论服务的GetCommentsByPostID方法获取指定文章的所有评论
	// 这个方法会查询数据库并返回评论列表（包含评论者信息）
	comments, err := c.service.GetCommentsByPostID(uint(postID))
	if err != nil {
		// 如果查询过程中发生错误，记录错误日志并返回500 Internal Server Error状态码
		middleware.Error("获取评论列表失败 (文章ID: " + postIDStr + "): " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论列表失败"})
		return // 提前返回，不再执行后续代码
	}

	// 格式化返回数据
	// 将评论列表转换为前端友好的格式
	var result []gin.H
	for _, comment := range comments {
		result = append(result, gin.H{
			"id":      comment.ID,       // 评论ID
			"content": comment.Content, // 评论内容
			"user": gin.H{ // 评论者信息
				"id":       comment.User.ID,       // 评论者ID
				"username": comment.User.Username, // 评论者用户名
			},
			"created_at": comment.CreatedAt.Format("2006-01-02 15:04:05"), // 格式化后的创建时间
		})
	}

	// 返回评论列表，状态码200 OK
	ctx.JSON(http.StatusOK, gin.H{
		"total":    len(result), // 评论总数
		"comments": result,     // 评论列表
	})
}

// DeleteComment 删除评论（仅评论作者可删除）
// @Summary 删除评论
// @Description 删除指定评论（仅评论作者可操作）
// @Tags comments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 403 {object} map[string]interface{} "无权操作"
// @Failure 404 {object} map[string]interface{} "评论不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/comments/{id} [delete]

// DeleteComment 处理删除评论请求
// 这是HTTP DELETE /api/comments/:id 路由的处理函数（需要认证）
// 参数:
//   ctx - Gin上下文对象，包含HTTP请求和响应的所有信息
func (c *CommentController) DeleteComment(ctx *gin.Context) {
	// 从Gin上下文中获取当前登录用户的ID
	userID, _ := ctx.Get("userID")
	
	// 从URL路径参数中获取评论ID字符串
	commentIDStr := ctx.Param("id")
	
	// 将字符串类型的评论ID转换为整数类型
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		// 如果转换失败，说明评论ID格式不正确，记录错误日志并返回400 Bad Request状态码
		middleware.Error("评论ID格式错误: " + commentIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论ID"})
		return // 提前返回，不再执行后续代码
	}

	// 验证评论是否存在且属于当前用户
	// 先获取评论详情，检查权限
	comment, err := c.service.GetCommentByID(uint(commentID))
	if err != nil {
		// 如果查询过程中发生错误
		if err == gorm.ErrRecordNotFound {
			// 如果错误是记录不存在，返回404 Not Found状态码
			ctx.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
			return
		}
		// 其他错误返回500 Internal Server Error状态码
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败"})
		return // 提前返回，不再执行后续代码
	}

	// 检查当前用户是否有权限删除这条评论
	// 只有评论的作者才能删除自己的评论
	if comment.UserID != userID.(uint) {
		// 如果当前用户不是评论作者，记录警告日志并返回403 Forbidden状态码
		middleware.Warn("无权删除评论 (用户ID: " + utils.UintToString(userID.(uint)) + ", 评论ID: " + commentIDStr + ")")
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权删除此评论"})
		return // 提前返回，不再执行后续代码
	}

	// 调用评论服务的DeleteComment方法删除评论
	// 这个方法会从数据库中删除评论记录
	if err := c.service.DeleteComment(uint(commentID)); err != nil {
		// 如果删除过程中发生错误，记录错误日志并返回500 Internal Server Error状态码
		middleware.Error("删除评论失败 (ID: " + commentIDStr + "): " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除评论失败"})
		return // 提前返回，不再执行后续代码
	}

	// 评论删除成功，记录信息日志并返回200 OK状态码
	middleware.Info("评论删除成功 (ID: " + commentIDStr + ")")
	ctx.JSON(http.StatusOK, gin.H{"message": "评论删除成功"}) // 返回成功消息
}