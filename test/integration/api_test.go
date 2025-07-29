package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"rakia-tech-test/internal/application/services"
	"rakia-tech-test/internal/infrastructure/repositories"
	"rakia-tech-test/internal/interfaces/rest"
)

type TestSuite struct {
	router *gin.Engine
	logger *logrus.Logger
}

func NewTestSuite() *TestSuite {
	gin.SetMode(gin.TestMode)

	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Suppress logs during tests

	postRepo := repositories.NewMemoryPostRepository()
	postService := services.NewPostService(postRepo, logger)
	postHandler := rest.NewPostHandler(postService, logger)

	r := rest.SetupRouter(postHandler, logger)

	return &TestSuite{
		router: r,
		logger: logger,
	}
}

func TestAPI_Health(t *testing.T) {
	suite := NewTestSuite()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestAPI_CreatePost(t *testing.T) {
	suite := NewTestSuite()

	testsCases := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectError    bool
	}{
		{
			name: "valid post creation",
			payload: map[string]interface{}{
				"title":   "Test Title",
				"content": "Test Content",
				"author":  "Test Author",
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "missing title",
			payload: map[string]interface{}{
				"content": "Test Content",
				"author":  "Test Author",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "empty title",
			payload: map[string]interface{}{
				"title":   "",
				"content": "Test Content",
				"author":  "Test Author",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "missing content",
			payload: map[string]interface{}{
				"title":  "Test Title",
				"author": "Test Author",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "missing author",
			payload: map[string]interface{}{
				"title":   "Test Title",
				"content": "Test Content",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tc := range testsCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonPayload, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/api/v1/posts", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			if tc.expectError {
				assert.Contains(t, response, "error")
			} else {
				assert.Contains(t, response, "id")
				assert.Equal(t, tc.payload["title"], response["title"])
				assert.Equal(t, tc.payload["content"], response["content"])
				assert.Equal(t, tc.payload["author"], response["author"])
				// ID should be a number (float64 in JSON)
				_, ok := response["id"].(float64)
				assert.True(t, ok, "ID should be a number")
			}
		})
	}
}

func TestAPI_GetPost(t *testing.T) {
	suite := NewTestSuite()

	// First create a post
	payload := map[string]interface{}{
		"title":   "Test Title",
		"content": "Test Content",
		"author":  "Test Author",
	}
	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/v1/posts", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var createdPost map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &createdPost)
	require.NoError(t, err)

	// Convert float64 to int for ID
	postIDFloat := createdPost["id"].(float64)
	postID := int(postIDFloat)

	testsCases := []struct {
		name           string
		postID         string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "get existing post",
			postID:         strconv.Itoa(postID),
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "get non-existent post",
			postID:         "999",
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:           "invalid post ID format",
			postID:         "invalid",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tc := range testsCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/posts/"+tc.postID, nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			if tc.expectError {
				assert.Contains(t, response, "error")
			} else {
				assert.Equal(t, float64(postID), response["id"])
				assert.Equal(t, payload["title"], response["title"])
				assert.Equal(t, payload["content"], response["content"])
				assert.Equal(t, payload["author"], response["author"])
			}
		})
	}
}

func TestAPI_GetAllPosts(t *testing.T) {
	suite := NewTestSuite()

	// Test empty list
	req, _ := http.NewRequest("GET", "/api/v1/posts", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(0), response["total"])
	assert.Empty(t, response["posts"])

	// Create a few posts
	posts := []map[string]interface{}{
		{"title": "Title 1", "content": "Content 1", "author": "Author 1"},
		{"title": "Title 2", "content": "Content 2", "author": "Author 2"},
		{"title": "Title 3", "content": "Content 3", "author": "Author 3"},
	}

	for _, post := range posts {
		jsonPayload, _ := json.Marshal(post)
		req, _ := http.NewRequest("POST", "/api/v1/posts", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
	}

	// Test populated list
	req, _ = http.NewRequest("GET", "/api/v1/posts", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(3), response["total"])
	postsArray := response["posts"].([]interface{})
	assert.Len(t, postsArray, 3)
}

func TestAPI_UpdatePost(t *testing.T) {
	suite := NewTestSuite()

	// First create a post
	payload := map[string]interface{}{
		"title":   "Original Title",
		"content": "Original Content",
		"author":  "Original Author",
	}
	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/v1/posts", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var createdPost map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &createdPost)
	require.NoError(t, err)

	// Convert float64 to int for ID
	postIDFloat := createdPost["id"].(float64)
	postID := int(postIDFloat)

	testsCases := []struct {
		name           string
		postID         string
		payload        map[string]interface{}
		expectedStatus int
		expectError    bool
	}{
		{
			name:   "valid update",
			postID: strconv.Itoa(postID),
			payload: map[string]interface{}{
				"title":   "Updated Title",
				"content": "Updated Content",
				"author":  "Updated Author",
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:   "update non-existent post",
			postID: "999",
			payload: map[string]interface{}{
				"title":   "Updated Title",
				"content": "Updated Content",
				"author":  "Updated Author",
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:   "invalid update data",
			postID: strconv.Itoa(postID),
			payload: map[string]interface{}{
				"title":   "",
				"content": "Updated Content",
				"author":  "Updated Author",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:   "invalid post ID format",
			postID: "invalid",
			payload: map[string]interface{}{
				"title":   "Updated Title",
				"content": "Updated Content",
				"author":  "Updated Author",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tc := range testsCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonPayload, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("PUT", "/api/v1/posts/"+tc.postID, bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			if tc.expectError {
				assert.Contains(t, response, "error")
			} else {
				assert.Equal(t, float64(postID), response["id"])
				assert.Equal(t, tc.payload["title"], response["title"])
				assert.Equal(t, tc.payload["content"], response["content"])
				assert.Equal(t, tc.payload["author"], response["author"])
			}
		})
	}
}

func TestAPI_DeletePost(t *testing.T) {
	suite := NewTestSuite()

	// First create a post
	payload := map[string]interface{}{
		"title":   "Test Title",
		"content": "Test Content",
		"author":  "Test Author",
	}
	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/v1/posts", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var createdPost map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &createdPost)
	require.NoError(t, err)

	// Convert float64 to int for ID
	postIDFloat := createdPost["id"].(float64)
	postID := int(postIDFloat)

	testsCases := []struct {
		name           string
		postID         string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "delete existing post",
			postID:         strconv.Itoa(postID),
			expectedStatus: http.StatusNoContent,
			expectError:    false,
		},
		{
			name:           "delete non-existent post",
			postID:         "999",
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:           "invalid post ID format",
			postID:         "invalid",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tc := range testsCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/api/v1/posts/"+tc.postID, nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			if tc.expectError && w.Code != http.StatusNoContent {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response, "error")
			}
		})
	}

	// Verify the post was actually deleted (only if the first test passed)
	if postID > 0 {
		req, _ = http.NewRequest("GET", "/api/v1/posts/"+strconv.Itoa(postID), nil)
		w = httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	}
}
