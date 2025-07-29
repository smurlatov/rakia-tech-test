package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"rakia-tech-test/internal/application/services"
	"rakia-tech-test/internal/domain/repositories"
	"rakia-tech-test/internal/interfaces/rest/dto"
)

type PostHandler struct {
	postService *services.PostService
	logger      *logrus.Logger
}

func NewPostHandler(postService *services.PostService, logger *logrus.Logger) *PostHandler {
	return &PostHandler{
		postService: postService,
		logger:      logger,
	}
}

// CreatePost handles POST /posts
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	post, err := h.postService.CreatePost(req.Title, req.Content, req.Author)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create post")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "creation_failed",
			Message: err.Error(),
		})
		return
	}

	response := dto.ToPostResponse(post)
	c.JSON(http.StatusCreated, response)
}

// GetPost handles GET /posts/:id
func (h *PostHandler) GetPost(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Post ID is required",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid post ID format",
		})
		return
	}

	post, err := h.postService.GetPostByID(id)
	if err != nil {
		if err == repositories.ErrPostNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Post not found",
			})
			return
		}

		h.logger.WithError(err).Error("Failed to get post")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve post",
		})
		return
	}

	response := dto.ToPostResponse(post)
	c.JSON(http.StatusOK, response)
}

// GetAllPosts handles GET /posts
func (h *PostHandler) GetAllPosts(c *gin.Context) {
	posts, err := h.postService.GetAllPosts()
	if err != nil {
		h.logger.WithError(err).Error("Failed to get posts")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve posts",
		})
		return
	}

	response := dto.ToPostsResponse(posts)
	c.JSON(http.StatusOK, response)
}

// UpdatePost handles PUT /posts/:id
func (h *PostHandler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Post ID is required",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid post ID format",
		})
		return
	}

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	post, err := h.postService.UpdatePost(id, req.Title, req.Content, req.Author)
	if err != nil {
		if err == repositories.ErrPostNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Post not found",
			})
			return
		}

		h.logger.WithError(err).Error("Failed to update post")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	response := dto.ToPostResponse(post)
	c.JSON(http.StatusOK, response)
}

// DeletePost handles DELETE /posts/:id
func (h *PostHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Post ID is required",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid post ID format",
		})
		return
	}

	err = h.postService.DeletePost(id)
	if err != nil {
		if err == repositories.ErrPostNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Post not found",
			})
			return
		}

		h.logger.WithError(err).Error("Failed to delete post")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete post",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
