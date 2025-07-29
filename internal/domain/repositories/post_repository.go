package repositories

import (
	"errors"
	"rakia-tech-test/internal/domain/entities"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrPostExists   = errors.New("post already exists")
)

type PostRepository interface {
	CreatePost(title, content, author string) (*entities.Post, error)

	Create(post *entities.Post) error

	GetByID(id int) (*entities.Post, error)

	GetAll() ([]*entities.Post, error)

	Update(id int, post *entities.Post) error

	Delete(id int) error

	Exists(id int) bool

	LoadData(posts []*entities.Post) error
}
