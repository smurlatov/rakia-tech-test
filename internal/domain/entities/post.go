package entities

import (
	"errors"
	"strings"
)

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

func NewPost(id int, title, content, author string) (*Post, error) {
	post := &Post{
		ID:      id,
		Title:   title,
		Content: content,
		Author:  author,
	}

	if err := post.Validate(); err != nil {
		return nil, err
	}

	return post, nil
}

func (p *Post) Update(title, content, author string) error {
	temp := &Post{
		ID:      p.ID,
		Title:   title,
		Content: content,
		Author:  author,
	}

	if err := temp.Validate(); err != nil {
		return err
	}

	p.Title = title
	p.Content = content
	p.Author = author

	return nil
}

func (p *Post) Validate() error {
	if strings.TrimSpace(p.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(p.Content) == "" {
		return errors.New("content is required")
	}
	if strings.TrimSpace(p.Author) == "" {
		return errors.New("author is required")
	}
	if len(p.Title) > 255 {
		return errors.New("title must be less than 255 characters")
	}
	return nil
}
