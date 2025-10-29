// controllers 包包含所有HTTP请求处理器的实现
// 这个文件专门处理用户相关的HTTP请求
package controllers

// 导入所需的包
import (
	"blogV2/internal/models"     // 数据模型包，包含数据库实体定义
	"blogV2/internal/services"   // 服务层包，包含业务逻辑实现
	"net/http"                   // HTTP协议包，包含HTTP状态码等常量

	"github.com/gin-gonic/gin"    // Gin Web框架
)

// UserController 用户控制器
// @Summary 用户管理接口
// @Description 提供用户注册、登录、查询等操作
// @Tags users
// @Accept json
// @Produce json
// @Router /api/v1/users [get]

// UserController 用户控制器结构体
// 负责处理所有与用户相关的HTTP请求
// 遵循MVC架构模式，作为视图层和模型层之间的桥梁
type UserController struct {
	service *services.UserService // 用户服务实例，用于处理业务逻辑
	                           // 通过依赖注入的方式注入，实现控制层与服务层的分离
}

// NewUserController 创建用户控制器实例
// 这是控制器的工厂函数，用于创建并初始化UserController实例
// 参数:
//   service - 用户服务实例，包含用户相关的业务逻辑
// 返回值:
//   *UserController - 初始化好的用户控制器实例
func NewUserController(service *services.UserService) *UserController {
	// 返回新创建的UserController实例
	// 将传入的服务实例赋值给控制器的service字段
	return &UserController{service: service}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账户
// @Tags users
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 201 {object} map[string]interface{} "注册成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/register [post]

// RegisterRequest 用户注册请求结构体
// 定义了客户端发送注册请求时需要提供的字段和验证规则
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"` // 用户名，必填字段，长度限制3-50个字符
	                                                               // binding标签用于Gin框架的请求验证
	Email    string `json:"email" binding:"required,email"`          // 邮箱地址，必填字段，需要符合邮箱格式验证
	                                                               // email验证器会检查字符串是否符合邮箱格式
	Password string `json:"password" binding:"required,min=6"`        // 密码，必填字段，最少6位字符
	                                                               // 密码会在服务层自动进行bcrypt加密
}

// Register 处理用户注册请求
// 这是HTTP POST /api/register 路由的处理函数
// 参数:
//   ctx - Gin上下文对象，包含HTTP请求和响应的所有信息
func (c *UserController) Register(ctx *gin.Context) {
	// 声明注册请求变量，用于存储解析后的请求数据
	var req RegisterRequest
	
	// 使用Gin框架的ShouldBindJSON方法将请求的JSON体绑定到RegisterRequest结构体
	// 这个方法会自动验证binding标签定义的规则
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 如果绑定或验证失败，返回400 Bad Request状态码
		// 包含具体的错误信息，帮助客户端了解验证失败的原因
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return // 提前返回，不再执行后续代码
	}

	// 创建用户模型实例，将请求数据转换为数据库实体
	// 这个实例将被传递给服务层进行业务处理
	user := &models.User{
		Username: req.Username, // 设置用户名，从请求中获取
		Email:    req.Email,    // 设置邮箱地址，从请求中获取
		Password: req.Password, // 设置密码，从请求中获取（服务层会自动加密）
	}

	// 调用用户服务的CreateUser方法创建新用户
	// 这个方法会处理密码加密、数据验证、数据库插入等业务逻辑
	if err := c.service.CreateUser(user); err != nil {
		// 如果创建用户过程中发生错误，返回500 Internal Server Error状态码
		// 这表示服务器端出现了未预期的错误
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return // 提前返回，不再执行后续代码
	}

	// 用户创建成功，返回201 Created状态码
	// 按照RESTful API最佳实践，创建资源成功应该返回201状态码
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully", // 成功消息，告知客户端注册成功
		"user": gin.H{ // 返回创建的用户信息（不包含敏感信息如密码）
			"id":       user.ID,       // 用户ID，由数据库自动生成
			"username": user.Username, // 用户名
			"email":    user.Email,   // 邮箱地址
		},
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取访问令牌
// @Tags users
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} map[string]interface{} "登录成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "用户名或密码错误"
// @Router /api/v1/login [post]

// LoginRequest 用户登录请求结构体
// 定义了客户端发送登录请求时需要提供的字段和验证规则
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名，必填字段
	Password string `json:"password" binding:"required"` // 密码，必填字段
}

// Login 处理用户登录请求
// 这是HTTP POST /api/login 路由的处理函数
// 参数:
//   ctx - Gin上下文对象，包含HTTP请求和响应的所有信息
func (c *UserController) Login(ctx *gin.Context) {
	// 声明登录请求变量，用于存储解析后的请求数据
	var req LoginRequest
	
	// 使用Gin框架的ShouldBindJSON方法将请求的JSON体绑定到LoginRequest结构体
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 如果绑定或验证失败，返回400 Bad Request状态码
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return // 提前返回，不再执行后续代码
	}

	// 调用用户服务的Login方法进行用户认证
	// 这个方法会验证用户名和密码，如果验证成功则生成JWT token
	// 返回值:
	//   user - 认证成功的用户信息
	//   token - 生成的JWT访问令牌
	//   err - 错误信息，如果认证失败则不为nil
	user, token, err := c.service.Login(req.Username, req.Password)
	if err != nil {
		// 登录失败，返回401 Unauthorized状态码
		// 这表示提供的凭据无效或用户不存在
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return // 提前返回，不再执行后续代码
	}

	// 登录成功，返回200 OK状态码和用户信息及访问令牌
	ctx.JSON(http.StatusOK, gin.H{
		"message": "login successful", // 成功消息
		"token":   token,              // JWT访问令牌，客户端需要保存这个token用于后续请求
		"user": gin.H{ // 用户基本信息
			"id":       user.ID,       // 用户ID
			"username": user.Username, // 用户名
			"email":    user.Email,   // 邮箱地址
		},
	})
}

// GetCurrentUser 获取当前登录用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "用户信息"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /api/v1/user [get]

// GetCurrentUser 获取当前登录用户信息
// 这是HTTP GET /api/user/me 路由的处理函数（需要认证）
// 参数:
//   ctx - Gin上下文对象，包含HTTP请求和响应的所有信息
func (c *UserController) GetCurrentUser(ctx *gin.Context) {
	// 从Gin上下文中获取当前登录用户的ID
	// 这个值由认证中间件在验证JWT token后设置
	userID, _ := ctx.Get("userID")
	
	// 调用用户服务的GetUserByID方法获取用户详细信息
	user, err := c.service.GetUserByID(userID.(uint))
	if err != nil {
		// 如果获取用户信息失败，返回500 Internal Server Error状态码
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return // 提前返回，不再执行后续代码
	}
	
	// 返回当前用户的详细信息，状态码200 OK
	ctx.JSON(http.StatusOK, gin.H{
		"id":         user.ID,        // 用户ID
		"username":   user.Username,  // 用户名
		"email":      user.Email,     // 邮箱地址
		"created_at": user.CreatedAt, // 账户创建时间
	})
}