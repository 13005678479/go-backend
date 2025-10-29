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

type CommentController struct {
	service *services.CommentService
}

func NewCommentController(service *services.CommentService) *CommentController {
	return &CommentController{service: service}
}

// CreateComment 创建评论
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=500"`
}

func (c *CommentController) CreateComment(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	postIDStr := ctx.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utils.Error("文章ID格式错误: " + postIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 验证文章是否存在
	post, err := c.service.GetPostByID(uint(postID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return
	}

	var req CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error("创建评论参数验证失败: " + err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return
	}

	comment := &models.Comment{
		Content: req.Content,
		PostID:  post.ID,
		UserID:  userID.(uint),
	}

	if err := c.service.CreateComment(comment); err != nil {
		utils.Error("创建评论失败 (文章ID: " + postIDStr + "): " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败"})
		return
	}

	utils.Info("评论创建成功 (ID: " + strconv.Itoa(int(comment.ID)) + ", 文章ID: " + postIDStr + ")")
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "评论创建成功",
		"comment": gin.H{
			"id":         comment.ID,
			"content":    comment.Content,
			"created_at": comment.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// ListComments 获取文章的所有评论
func (c *CommentController) ListComments(ctx *gin.Context) {
	postIDStr := ctx.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		utils.Error("文章ID格式错误: " + postIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 验证文章是否存在
	_, err = c.service.GetPostByID(uint(postID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return
	}

	comments, err := c.service.GetCommentsByPostID(uint(postID))
	if err != nil {
		utils.Error("获取评论列表失败 (文章ID: " + postIDStr + "): " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论列表失败"})
		return
	}

	// 格式化返回数据
	var result []gin.H
	for _, comment := range comments {
		result = append(result, gin.H{
			"id":      comment.ID,
			"content": comment.Content,
			"user": gin.H{
				"id":       comment.User.ID,
				"username": comment.User.Username,
			},
			"created_at": comment.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total":    len(result),
		"comments": result,
	})
}

// DeleteComment 删除评论（仅评论作者可删除）
func (c *CommentController) DeleteComment(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	commentIDStr := ctx.Param("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		utils.Error("评论ID格式错误: " + commentIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论ID"})
		return
	}

	// 验证评论是否存在且属于当前用户
	comment, err := c.service.GetCommentByID(uint(commentID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败"})
		return
	}

	if comment.UserID != userID.(uint) {
		utils.Warn("无权删除评论 (用户ID: " + utils.UintToString(userID.(uint)) + ", 评论ID: " + commentIDStr + ")")
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权删除此评论"})
		return
	}

	if err := c.service.DeleteComment(uint(commentID)); err != nil {
		utils.Error("删除评论失败 (ID: " + commentIDStr + "): " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除评论失败"})
		return
	}

	utils.Info("评论删除成功 (ID: " + commentIDStr + ")")
	ctx.JSON(http.StatusOK, gin.H{"message": "评论删除成功"})
}
