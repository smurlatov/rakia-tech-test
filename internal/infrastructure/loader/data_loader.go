package loader

import (
	"encoding/json"
	"os"
	"rakia-tech-test/internal/domain/entities"
	"rakia-tech-test/internal/domain/repositories"

	"github.com/sirupsen/logrus"
)

type PostData struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

type BlogData struct {
	Posts []PostData `json:"posts"`
}

type DataLoader struct {
	postRepo repositories.PostRepository
	logger   *logrus.Logger
}

func NewDataLoader(postRepo repositories.PostRepository, logger *logrus.Logger) *DataLoader {
	return &DataLoader{
		postRepo: postRepo,
		logger:   logger,
	}
}

func (dl *DataLoader) LoadFromFile(filename string) error {
	dl.logger.WithField("filename", filename).Info("Loading blog data from file")

	data, err := os.ReadFile(filename)
	if err != nil {
		dl.logger.WithError(err).Error("Failed to read data file")
		return err
	}

	var blogData BlogData
	if err := json.Unmarshal(data, &blogData); err != nil {
		dl.logger.WithError(err).Error("Failed to parse JSON data")
		return err
	}

	posts := make([]*entities.Post, len(blogData.Posts))
	for i, postData := range blogData.Posts {
		post, err := entities.NewPost(postData.ID, postData.Title, postData.Content, postData.Author)
		if err != nil {
			dl.logger.WithError(err).WithField("post_id", postData.ID).Error("Failed to create post entity")
			return err
		}
		posts[i] = post
	}

	if err := dl.postRepo.LoadData(posts); err != nil {
		dl.logger.WithError(err).Error("Failed to load data into repository")
		return err
	}

	dl.logger.WithField("count", len(posts)).Info("Successfully loaded blog posts")
	return nil
}
