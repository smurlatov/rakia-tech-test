package dto

import (
	"rakia-tech-test/internal/domain/entities"
)

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,max=255"`
	Content string `json:"content" binding:"required"`
	Author  string `json:"author" binding:"required"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required,max=255"`
	Content string `json:"content" binding:"required"`
	Author  string `json:"author" binding:"required"`
}

type PostResponse struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type PostsResponse struct {
	Posts []PostResponse `json:"posts"`
	Total int            `json:"total"`
}

func ToPostResponse(post *entities.Post) PostResponse {
	return PostResponse{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Author:  post.Author,
	}
}

func ToPostsResponse(posts []*entities.Post) PostsResponse {
	responses := make([]PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = ToPostResponse(post)
	}

	return PostsResponse{
		Posts: responses,
		Total: len(posts),
	}
}
