package controllers

import (
	"blogV2/internal/models"
	"blogV2/internal/services"
	"blogV2/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	service *services.PostService
}

func NewPostController(service *services.PostService) *PostController {
	return &PostController{service: service}
}

// CreatePost 创建文章
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
