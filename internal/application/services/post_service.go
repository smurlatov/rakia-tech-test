package services

import (
	"rakia-tech-test/internal/domain/entities"
	"rakia-tech-test/internal/domain/repositories"

	"github.com/sirupsen/logrus"
)

type PostService struct {
	postRepo repositories.PostRepository
	logger   *logrus.Logger
}

func NewPostService(postRepo repositories.PostRepository, logger *logrus.Logger) *PostService {
	return &PostService{
		postRepo: postRepo,
		logger:   logger,
	}
}

func (s *PostService) CreatePost(title, content, author string) (*entities.Post, error) {
	s.logger.WithFields(logrus.Fields{
		"title":  title,
		"author": author,
	}).Info("Creating new post")

	post, err := s.postRepo.CreatePost(title, content, author)
	if err != nil {
		return nil, err
	}

	s.logger.WithField("post_id", post.ID).Info("Post created successfully")
	return post, nil
}

func (s *PostService) GetPostByID(id int) (*entities.Post, error) {
	s.logger.WithField("post_id", id).Debug("Retrieving post by ID")

	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetAllPosts() ([]*entities.Post, error) {
	s.logger.Debug("Retrieving all posts")

	posts, err := s.postRepo.GetAll()
	if err != nil {
		return nil, err
	}

	s.logger.WithField("count", len(posts)).Debug("Retrieved posts")
	return posts, nil
}

func (s *PostService) UpdatePost(id int, title, content, author string) (*entities.Post, error) {
	s.logger.WithFields(logrus.Fields{
		"post_id": id,
		"title":   title,
		"author":  author,
	}).Info("Updating post")

	existingPost, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := existingPost.Update(title, content, author); err != nil {
		return nil, err
	}

	if err := s.postRepo.Update(id, existingPost); err != nil {
		return nil, err
	}

	s.logger.WithField("post_id", id).Info("Post updated successfully")
	return existingPost, nil
}

func (s *PostService) DeletePost(id int) error {
	s.logger.WithField("post_id", id).Info("Deleting post")

	if err := s.postRepo.Delete(id); err != nil {
		return err
	}

	s.logger.WithField("post_id", id).Info("Post deleted successfully")
	return nil
}
