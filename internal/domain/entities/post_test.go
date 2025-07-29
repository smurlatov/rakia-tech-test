package entities

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test fixtures and data providers
type newPostTestCase struct {
	name          string
	id            int
	title         string
	content       string
	author        string
	wantError     bool
	expectedError string
}

type updateTestCase struct {
	name          string
	title         string
	content       string
	author        string
	wantError     bool
	expectedError string
}

type validateTestCase struct {
	name          string
	post          Post
	wantErr       bool
	expectedError string
}

// Test data providers
func getNewPostTestCases() []newPostTestCase {
	longTitle256 := string(make([]byte, 256))
	validTitle255 := string(make([]byte, 255))

	return []newPostTestCase{
		{
			name:      "valid post",
			id:        1,
			title:     "Test Title",
			content:   "Test Content",
			author:    "Test Author",
			wantError: false,
		},
		{
			name:          "empty title",
			id:            2,
			title:         "",
			content:       "Test Content",
			author:        "Test Author",
			wantError:     true,
			expectedError: "title is required",
		},
		{
			name:          "whitespace only title",
			id:            3,
			title:         "   ",
			content:       "Test Content",
			author:        "Test Author",
			wantError:     true,
			expectedError: "title is required",
		},
		{
			name:          "empty content",
			id:            4,
			title:         "Test Title",
			content:       "",
			author:        "Test Author",
			wantError:     true,
			expectedError: "content is required",
		},
		{
			name:          "whitespace only content",
			id:            5,
			title:         "Test Title",
			content:       "   ",
			author:        "Test Author",
			wantError:     true,
			expectedError: "content is required",
		},
		{
			name:          "empty author",
			id:            6,
			title:         "Test Title",
			content:       "Test Content",
			author:        "",
			wantError:     true,
			expectedError: "author is required",
		},
		{
			name:          "whitespace only author",
			id:            7,
			title:         "Test Title",
			content:       "Test Content",
			author:        "   ",
			wantError:     true,
			expectedError: "author is required",
		},
		{
			name:          "title too long",
			id:            8,
			title:         longTitle256,
			content:       "Test Content",
			author:        "Test Author",
			wantError:     true,
			expectedError: "title must be less than 255 characters",
		},
		{
			name:      "title exactly 255 characters (valid)",
			id:        9,
			title:     validTitle255,
			content:   "Test Content",
			author:    "Test Author",
			wantError: false,
		},
		{
			name:          "title 300 characters (invalid)",
			id:            10,
			title:         strings.Repeat("a", 300),
			content:       "Test Content",
			author:        "Test Author",
			wantError:     true,
			expectedError: "title must be less than 255 characters",
		},
	}
}

func getUpdateTestCases() []updateTestCase {
	return []updateTestCase{
		{
			name:      "valid update",
			title:     "Updated Title",
			content:   "Updated Content",
			author:    "Updated Author",
			wantError: false,
		},
		{
			name:          "empty title update",
			title:         "",
			content:       "Updated Content",
			author:        "Updated Author",
			wantError:     true,
			expectedError: "title is required",
		},
		{
			name:          "whitespace only title update",
			title:         "   ",
			content:       "Updated Content",
			author:        "Updated Author",
			wantError:     true,
			expectedError: "title is required",
		},
		{
			name:          "empty content update",
			title:         "Updated Title",
			content:       "",
			author:        "Updated Author",
			wantError:     true,
			expectedError: "content is required",
		},
		{
			name:          "empty author update",
			title:         "Updated Title",
			content:       "Updated Content",
			author:        "",
			wantError:     true,
			expectedError: "author is required",
		},
		{
			name:          "title too long update",
			title:         strings.Repeat("x", 256),
			content:       "Updated Content",
			author:        "Updated Author",
			wantError:     true,
			expectedError: "title must be less than 255 characters",
		},
	}
}

func getValidateTestCases() []validateTestCase {
	return []validateTestCase{
		{
			name: "valid post",
			post: Post{
				ID:      1,
				Title:   "Valid Title",
				Content: "Valid Content",
				Author:  "Valid Author",
			},
			wantErr: false,
		},
		{
			name: "empty title validation",
			post: Post{
				ID:      2,
				Title:   "",
				Content: "Valid Content",
				Author:  "Valid Author",
			},
			wantErr:       true,
			expectedError: "title is required",
		},
		{
			name: "whitespace only title validation",
			post: Post{
				ID:      3,
				Title:   "   ",
				Content: "Valid Content",
				Author:  "Valid Author",
			},
			wantErr:       true,
			expectedError: "title is required",
		},
		{
			name: "tab and newline title validation",
			post: Post{
				ID:      4,
				Title:   "\t\n\r",
				Content: "Valid Content",
				Author:  "Valid Author",
			},
			wantErr:       true,
			expectedError: "title is required",
		},
		{
			name: "empty content validation",
			post: Post{
				ID:      5,
				Title:   "Valid Title",
				Content: "",
				Author:  "Valid Author",
			},
			wantErr:       true,
			expectedError: "content is required",
		},
		{
			name: "whitespace only content validation",
			post: Post{
				ID:      6,
				Title:   "Valid Title",
				Content: "   ",
				Author:  "Valid Author",
			},
			wantErr:       true,
			expectedError: "content is required",
		},
		{
			name: "empty author validation",
			post: Post{
				ID:      7,
				Title:   "Valid Title",
				Content: "Valid Content",
				Author:  "",
			},
			wantErr:       true,
			expectedError: "author is required",
		},
		{
			name: "whitespace only author validation",
			post: Post{
				ID:      8,
				Title:   "Valid Title",
				Content: "Valid Content",
				Author:  "   ",
			},
			wantErr:       true,
			expectedError: "author is required",
		},
		{
			name: "title too long validation",
			post: Post{
				ID:      9,
				Title:   strings.Repeat("a", 256),
				Content: "Valid Content",
				Author:  "Valid Author",
			},
			wantErr:       true,
			expectedError: "title must be less than 255 characters",
		},
		{
			name: "title exactly 255 characters (valid)",
			post: Post{
				ID:      10,
				Title:   strings.Repeat("b", 255),
				Content: "Valid Content",
				Author:  "Valid Author",
			},
			wantErr: false,
		},
	}
}

// Actual test functions
func TestNewPost(t *testing.T) {
	testCases := getNewPostTestCases()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			post, err := NewPost(tt.id, tt.title, tt.content, tt.author)

			if tt.wantError {
				require.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, post)
			} else {
				require.NoError(t, err)
				require.NotNil(t, post)
				assert.Equal(t, tt.id, post.ID)
				assert.Equal(t, tt.title, post.Title)
				assert.Equal(t, tt.content, post.Content)
				assert.Equal(t, tt.author, post.Author)
			}
		})
	}
}

func TestPost_Update(t *testing.T) {
	testCases := getUpdateTestCases()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			post, err := NewPost(1, "Original Title", "Original Content", "Original Author")
			require.NoError(t, err)

			err = post.Update(tt.title, tt.content, tt.author)

			if tt.wantError {
				require.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				assert.Equal(t, "Original Title", post.Title)
				assert.Equal(t, "Original Content", post.Content)
				assert.Equal(t, "Original Author", post.Author)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.title, post.Title)
				assert.Equal(t, tt.content, post.Content)
				assert.Equal(t, tt.author, post.Author)
			}
		})
	}
}

func TestPost_Validate(t *testing.T) {
	testCases := getValidateTestCases()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.post.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
