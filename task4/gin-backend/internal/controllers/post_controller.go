package controllers

import (
	"blogV2/internal/models"
	"blogV2/internal/services"
	"blogV2/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PostController 文章控制器
// @Summary 文章管理接口
// @Description 提供文章的创建、查询、更新、删除等操作
// @Tags posts
// @Accept json
// @Produce json
// @Router /api/v1/posts [get]

type PostController struct {
	service *services.PostService
}

func NewPostController(service *services.PostService) *PostController {
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
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=200"`
	Content string `json:"content" binding:"required,min=10"`
}

func (c *PostController) CreatePost(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	var req CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error("创建文章参数验证失败: " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return
	}

	post := &models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}

	if err := c.service.CreatePost(post); err != nil {
		utils.Error("创建文章失败 (用户ID: " + utils.UintToString(userID.(uint)) + "): " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}

	utils.Info("文章创建成功 (ID: " + strconv.Itoa(int(post.ID)) + ", 用户: " + utils.UintToString(userID.(uint)) + ")")
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "文章创建成功",
		"post": gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"content":    post.Content,
			"created_at": post.CreatedAt.Format("2006-01-02 15:04:05"),
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
func (c *PostController) GetPost(ctx *gin.Context) {
	postIDStr := ctx.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utils.Error("文章ID格式错误: " + postIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	post, err := c.service.GetPostByID(uint(postID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"post": gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"content":    post.Content,
			"user_id":    post.UserID,
			"username":   post.User.Username,
			"created_at": post.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at": post.UpdatedAt.Format("2006-01-02 15:04:05"),
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
func (c *PostController) UpdatePost(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	postIDStr := ctx.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utils.Error("文章ID格式错误: " + postIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var req CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error("更新文章参数验证失败: " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return
	}

	// 验证文章是否存在且属于当前用户
	post, err := c.service.GetPostByID(uint(postID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return
	}

	if post.UserID != userID.(uint) {
		utils.Warn("无权更新文章 (用户ID: " + utils.UintToString(userID.(uint)) + ", 文章ID: " + postIDStr + ")")
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权更新此文章"})
		return
	}

	// 更新文章
	post.Title = req.Title
	post.Content = req.Content
	
	if err := c.service.UpdatePost(post); err != nil {
		utils.Error("更新文章失败 (ID: " + postIDStr + "): " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}

	utils.Info("文章更新成功 (ID: " + postIDStr + ")")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "文章更新成功",
		"post": gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"content":    post.Content,
			"updated_at": post.UpdatedAt.Format("2006-01-02 15:04:05"),
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
func (c *PostController) DeletePost(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	postIDStr := ctx.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utils.Error("文章ID格式错误: " + postIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 验证文章是否存在且属于当前用户
	post, err := c.service.GetPostByID(uint(postID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return
	}

	if post.UserID != userID.(uint) {
		utils.Warn("无权删除文章 (用户ID: " + utils.UintToString(userID.(uint)) + ", 文章ID: " + postIDStr + ")")
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权删除此文章"})
		return
	}

	if err := c.service.DeletePost(uint(postID)); err != nil {
		utils.Error("删除文章失败 (ID: " + postIDStr + "): " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败"})
		return
	}

	utils.Info("文章删除成功 (ID: " + postIDStr + ")")
	ctx.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
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
func (c *PostController) ListPosts(ctx *gin.Context) {
	posts, err := c.service.ListAllPosts()
	if err != nil {
		utils.Error("获取文章列表失败: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}

	// 格式化返回数据
	var result []gin.H
	for _, post := range posts {
		result = append(result, gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"content":    post.Content,
			"user_id":    post.UserID,
			"username":   post.User.Username,
			"created_at": post.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at": post.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": len(result),
		"posts": result,
	})
}
