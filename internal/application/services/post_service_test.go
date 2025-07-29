package services

import (
	"errors"
	"rakia-tech-test/internal/domain/entities"
	"rakia-tech-test/internal/domain/repositories"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Create(post *entities.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockPostRepository) GetByID(id int) (*entities.Post, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Post), args.Error(1)
}

func (m *MockPostRepository) GetAll() ([]*entities.Post, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Post), args.Error(1)
}

func (m *MockPostRepository) Update(id int, post *entities.Post) error {
	args := m.Called(id, post)
	return args.Error(0)
}

func (m *MockPostRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPostRepository) Exists(id int) bool {
	args := m.Called(id)
	return args.Bool(0)
}

func (m *MockPostRepository) CreatePost(title, content, author string) (*entities.Post, error) {
	args := m.Called(title, content, author)
	if len(args) >= 2 && args.Get(0) != nil {
		return args.Get(0).(*entities.Post), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPostRepository) LoadData(posts []*entities.Post) error {
	args := m.Called(posts)
	return args.Error(0)
}

// Test case types
type createPostTestCase struct {
	name      string
	title     string
	content   string
	author    string
	mockSetup func(*MockPostRepository) *entities.Post
	wantError bool
}

type getPostByIDTestCase struct {
	name      string
	id        int
	mockSetup func(*MockPostRepository, *entities.Post)
	wantError bool
	wantPost  *entities.Post
}

type getAllPostsTestCase struct {
	name      string
	mockSetup func(*MockPostRepository, []*entities.Post)
	wantError bool
	wantCount int
}

type updatePostTestCase struct {
	name      string
	id        int
	title     string
	content   string
	author    string
	mockSetup func(*MockPostRepository, *entities.Post)
	wantError bool
}

type deletePostTestCase struct {
	name      string
	id        int
	mockSetup func(*MockPostRepository)
	wantError bool
}

// Test data providers
func getCreatePostTestCases() []createPostTestCase {
	return []createPostTestCase{
		{
			name:    "successful creation",
			title:   "Test Title",
			content: "Test Content",
			author:  "Test Author",
			mockSetup: func(mockRepo *MockPostRepository) *entities.Post {
				post, _ := entities.NewPost(1, "Test Title", "Test Content", "Test Author")
				mockRepo.On("CreatePost", "Test Title", "Test Content", "Test Author").Return(post, nil).Once()
				return post
			},
			wantError: false,
		},
		{
			name:    "validation error",
			title:   "",
			content: "Test Content",
			author:  "Test Author",
			mockSetup: func(mockRepo *MockPostRepository) *entities.Post {
				mockRepo.On("CreatePost", "", "Test Content", "Test Author").Return(nil, errors.New("title is required")).Once()
				return nil
			},
			wantError: true,
		},
		{
			name:    "repository error",
			title:   "Test Title",
			content: "Test Content",
			author:  "Test Author",
			mockSetup: func(mockRepo *MockPostRepository) *entities.Post {
				mockRepo.On("CreatePost", "Test Title", "Test Content", "Test Author").Return(nil, errors.New("repository error")).Once()
				return nil
			},
			wantError: true,
		},
	}
}

func getGetPostByIDTestCases() []getPostByIDTestCase {
	testPost, _ := entities.NewPost(1, "Test Title", "Test Content", "Test Author")

	return []getPostByIDTestCase{
		{
			name: "successful retrieval",
			id:   testPost.ID,
			mockSetup: func(mockRepo *MockPostRepository, post *entities.Post) {
				mockRepo.On("GetByID", post.ID).Return(post, nil).Once()
			},
			wantError: false,
			wantPost:  testPost,
		},
		{
			name: "post not found",
			id:   999,
			mockSetup: func(mockRepo *MockPostRepository, post *entities.Post) {
				mockRepo.On("GetByID", 999).Return(nil, repositories.ErrPostNotFound).Once()
			},
			wantError: true,
			wantPost:  nil,
		},
	}
}

func getGetAllPostsTestCases() []getAllPostsTestCase {
	return []getAllPostsTestCase{
		{
			name: "successful retrieval",
			mockSetup: func(mockRepo *MockPostRepository, posts []*entities.Post) {
				mockRepo.On("GetAll").Return(posts, nil).Once()
			},
			wantError: false,
			wantCount: 2,
		},
		{
			name: "repository error",
			mockSetup: func(mockRepo *MockPostRepository, posts []*entities.Post) {
				mockRepo.On("GetAll").Return(nil, errors.New("repository error")).Once()
			},
			wantError: true,
			wantCount: 0,
		},
	}
}

func getUpdatePostTestCases() []updatePostTestCase {
	return []updatePostTestCase{
		{
			name:    "successful update",
			id:      1,
			title:   "Updated Title",
			content: "Updated Content",
			author:  "Updated Author",
			mockSetup: func(mockRepo *MockPostRepository, post *entities.Post) {
				mockRepo.On("GetByID", post.ID).Return(post, nil).Once()
				mockRepo.On("Update", post.ID, mock.AnythingOfType("*entities.Post")).Return(nil).Once()
			},
			wantError: false,
		},
		{
			name:    "post not found",
			id:      999,
			title:   "Updated Title",
			content: "Updated Content",
			author:  "Updated Author",
			mockSetup: func(mockRepo *MockPostRepository, post *entities.Post) {
				mockRepo.On("GetByID", 999).Return(nil, repositories.ErrPostNotFound).Once()
			},
			wantError: true,
		},
		{
			name:    "validation error",
			id:      1,
			title:   "", // Invalid title
			content: "Updated Content",
			author:  "Updated Author",
			mockSetup: func(mockRepo *MockPostRepository, post *entities.Post) {
				mockRepo.On("GetByID", post.ID).Return(post, nil).Once()
			},
			wantError: true,
		},
	}
}

func getDeletePostTestCases() []deletePostTestCase {
	return []deletePostTestCase{
		{
			name: "successful deletion",
			id:   1,
			mockSetup: func(mockRepo *MockPostRepository) {
				mockRepo.On("Delete", 1).Return(nil).Once()
			},
			wantError: false,
		},
		{
			name: "post not found",
			id:   999,
			mockSetup: func(mockRepo *MockPostRepository) {
				mockRepo.On("Delete", 999).Return(repositories.ErrPostNotFound).Once()
			},
			wantError: true,
		},
	}
}

// Tests
func TestPostService_CreatePost(t *testing.T) {
	mockRepo := new(MockPostRepository)
	logger := logrus.New()
	service := NewPostService(mockRepo, logger)

	testCases := getCreatePostTestCases()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			expectedPost := tt.mockSetup(mockRepo)

			post, err := service.CreatePost(tt.title, tt.content, tt.author)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, post)
			} else {
				require.NoError(t, err)
				require.NotNil(t, post)
				assert.Equal(t, expectedPost.ID, post.ID)
				assert.Equal(t, tt.title, post.Title)
				assert.Equal(t, tt.content, post.Content)
				assert.Equal(t, tt.author, post.Author)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPostService_GetPostByID(t *testing.T) {
	mockRepo := new(MockPostRepository)
	logger := logrus.New()
	service := NewPostService(mockRepo, logger)

	testCases := getGetPostByIDTestCases()
	testPost, _ := entities.NewPost(1, "Test Title", "Test Content", "Test Author")

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup(mockRepo, testPost)

			result, err := service.GetPostByID(tt.id)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.wantPost.ID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPostService_GetAllPosts(t *testing.T) {
	mockRepo := new(MockPostRepository)
	logger := logrus.New()
	service := NewPostService(mockRepo, logger)

	post1, _ := entities.NewPost(1, "Title 1", "Content 1", "Author 1")
	post2, _ := entities.NewPost(2, "Title 2", "Content 2", "Author 2")
	posts := []*entities.Post{post1, post2}

	testCases := getGetAllPostsTestCases()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup(mockRepo, posts)

			result, err := service.GetAllPosts()

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPostService_UpdatePost(t *testing.T) {
	mockRepo := new(MockPostRepository)
	logger := logrus.New()
	service := NewPostService(mockRepo, logger)

	existingPost, _ := entities.NewPost(1, "Original Title", "Original Content", "Original Author")
	testCases := getUpdatePostTestCases()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup(mockRepo, existingPost)

			result, err := service.UpdatePost(tt.id, tt.title, tt.content, tt.author)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.title, result.Title)
				assert.Equal(t, tt.content, result.Content)
				assert.Equal(t, tt.author, result.Author)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPostService_DeletePost(t *testing.T) {
	mockRepo := new(MockPostRepository)
	logger := logrus.New()
	service := NewPostService(mockRepo, logger)

	testCases := getDeletePostTestCases()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockSetup(mockRepo)

			err := service.DeletePost(tt.id)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
