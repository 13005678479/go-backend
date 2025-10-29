package controllers

import (
	"blogV2/internal/models"
	"blogV2/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
// @Summary 用户管理接口
// @Description 提供用户注册、登录、查询等操作
// @Tags users
// @Accept json
// @Produce json
// @Router /api/v1/users [get]

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
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
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (c *UserController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := c.service.CreateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
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
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *UserController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := c.service.Login(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
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
func (c *UserController) GetCurrentUser(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	user, err := c.service.GetUserByID(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		// "post_count": user.PostCount,
		"created_at": user.CreatedAt,
	})
}
