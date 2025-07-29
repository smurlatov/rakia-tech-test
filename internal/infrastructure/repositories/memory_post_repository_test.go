package repositories

import (
	"fmt"
	"rakia-tech-test/internal/domain/entities"
	"rakia-tech-test/internal/domain/repositories"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryPostRepository_Create(t *testing.T) {
	repo := NewMemoryPostRepository()

	post, err := entities.NewPost(1, "Test Title", "Test Content", "Test Author")
	require.NoError(t, err)

	err = repo.Create(post)
	assert.NoError(t, err)

	retrievedPost, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, post.ID, retrievedPost.ID)
	assert.Equal(t, post.Title, retrievedPost.Title)
	assert.Equal(t, post.Content, retrievedPost.Content)
	assert.Equal(t, post.Author, retrievedPost.Author)

	err = repo.Create(post)
	assert.ErrorIs(t, err, repositories.ErrPostExists)
}

func TestMemoryPostRepository_CreatePost(t *testing.T) {
	repo := NewMemoryPostRepository()

	post, err := repo.CreatePost("Test Title", "Test Content", "Test Author")
	require.NoError(t, err)
	require.NotNil(t, post)

	assert.Equal(t, 1, post.ID)
	assert.Equal(t, "Test Title", post.Title)
	assert.Equal(t, "Test Content", post.Content)
	assert.Equal(t, "Test Author", post.Author)

	retrievedPost, err := repo.GetByID(post.ID)
	require.NoError(t, err)
	assert.Equal(t, post.ID, retrievedPost.ID)
	assert.Equal(t, post.Title, retrievedPost.Title)
	assert.Equal(t, post.Content, retrievedPost.Content)
	assert.Equal(t, post.Author, retrievedPost.Author)

	post2, err := repo.CreatePost("Second Title", "Second Content", "Second Author")
	require.NoError(t, err)
	assert.Equal(t, 2, post2.ID)

	_, err = repo.CreatePost("", "Content", "Author")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "title is required")

	var wg sync.WaitGroup
	var posts []*entities.Post
	var mu sync.Mutex

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			post, err := repo.CreatePost(
				fmt.Sprintf("Title %d", i),
				fmt.Sprintf("Content %d", i),
				fmt.Sprintf("Author %d", i),
			)
			if err == nil {
				mu.Lock()
				posts = append(posts, post)
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	assert.Len(t, posts, 10)
	idMap := make(map[int]bool)
	for _, p := range posts {
		assert.False(t, idMap[p.ID], "Duplicate ID found: %d", p.ID)
		idMap[p.ID] = true
	}
}

func TestMemoryPostRepository_GetByID(t *testing.T) {
	repo := NewMemoryPostRepository()

	post, err := entities.NewPost(1, "Test Title", "Test Content", "Test Author")
	require.NoError(t, err)

	_, err = repo.GetByID(post.ID)
	assert.Equal(t, repositories.ErrPostNotFound, err)

	err = repo.Create(post)
	require.NoError(t, err)

	retrievedPost, err := repo.GetByID(post.ID)
	require.NoError(t, err)
	assert.Equal(t, post.ID, retrievedPost.ID)
	assert.Equal(t, post.Title, retrievedPost.Title)
	assert.Equal(t, post.Content, retrievedPost.Content)
	assert.Equal(t, post.Author, retrievedPost.Author)

	retrievedPost.Title = "Modified Title"
	originalPost, err := repo.GetByID(post.ID)
	require.NoError(t, err)
	assert.Equal(t, post.Title, originalPost.Title)
}

func TestMemoryPostRepository_GetAll(t *testing.T) {
	repo := NewMemoryPostRepository()

	posts, err := repo.GetAll()
	require.NoError(t, err)
	assert.Empty(t, posts)

	// Create posts with IDs out of order to test sorting
	post3, err := entities.NewPost(3, "Title 3", "Content 3", "Author 3")
	require.NoError(t, err)
	post1, err := entities.NewPost(1, "Title 1", "Content 1", "Author 1")
	require.NoError(t, err)
	post2, err := entities.NewPost(2, "Title 2", "Content 2", "Author 2")
	require.NoError(t, err)

	// Add posts in random order
	err = repo.Create(post3)
	require.NoError(t, err)
	err = repo.Create(post1)
	require.NoError(t, err)
	err = repo.Create(post2)
	require.NoError(t, err)

	posts, err = repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, posts, 3)

	// Verify posts are sorted by ID in ascending order
	assert.Equal(t, 1, posts[0].ID)
	assert.Equal(t, 2, posts[1].ID)
	assert.Equal(t, 3, posts[2].ID)

	// Test defensive copying
	posts[0].Title = "Modified Title"
	originalPosts, err := repo.GetAll()
	require.NoError(t, err)
	assert.NotEqual(t, "Modified Title", originalPosts[0].Title)
	assert.NotEqual(t, "Modified Title", originalPosts[1].Title)
	assert.NotEqual(t, "Modified Title", originalPosts[2].Title)
}

func TestMemoryPostRepository_Update(t *testing.T) {
	repo := NewMemoryPostRepository()

	post, err := entities.NewPost(1, "Test Title", "Test Content", "Test Author")
	require.NoError(t, err)

	err = repo.Update(post.ID, post)
	assert.Equal(t, repositories.ErrPostNotFound, err)

	err = repo.Create(post)
	require.NoError(t, err)

	updatedPost := *post
	updatedPost.Title = "Updated Title"
	err = repo.Update(post.ID, &updatedPost)
	require.NoError(t, err)

	retrievedPost, err := repo.GetByID(post.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", retrievedPost.Title)
	assert.Equal(t, post.ID, retrievedPost.ID) // ID should remain the same
}

func TestMemoryPostRepository_Delete(t *testing.T) {
	repo := NewMemoryPostRepository()

	post, err := entities.NewPost(1, "Test Title", "Test Content", "Test Author")
	require.NoError(t, err)

	err = repo.Delete(post.ID)
	assert.Equal(t, repositories.ErrPostNotFound, err)

	err = repo.Create(post)
	require.NoError(t, err)

	err = repo.Delete(post.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(post.ID)
	assert.Equal(t, repositories.ErrPostNotFound, err)
}

func TestMemoryPostRepository_Exists(t *testing.T) {
	repo := NewMemoryPostRepository()

	post, err := entities.NewPost(1, "Test Title", "Test Content", "Test Author")
	require.NoError(t, err)

	assert.False(t, repo.Exists(post.ID))

	err = repo.Create(post)
	require.NoError(t, err)

	assert.True(t, repo.Exists(post.ID))
}

func TestMemoryPostRepository_LoadData(t *testing.T) {
	repo := NewMemoryPostRepository()

	posts := []*entities.Post{
		{ID: 1, Title: "Post 1", Content: "Content 1", Author: "Author 1"},
		{ID: 3, Title: "Post 3", Content: "Content 3", Author: "Author 3"},
		{ID: 2, Title: "Post 2", Content: "Content 2", Author: "Author 2"},
	}

	err := repo.LoadData(posts)
	assert.NoError(t, err)

	allPosts, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, allPosts, 3)

	// Verify nextID was updated correctly (should be max ID + 1)
	// Create a new post to check nextID
	newPost, err := repo.CreatePost("New Title", "New Content", "New Author")
	require.NoError(t, err)
	assert.Equal(t, 4, newPost.ID) // Should be 4 (max loaded ID 3 + 1)
}

func TestMemoryPostRepository_ConcurrentAccess(t *testing.T) {
	repo := NewMemoryPostRepository()
	var wg sync.WaitGroup

	// Create posts concurrently
	numGoroutines := 10
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()

			post, err := entities.NewPost(
				index+1, // Use index+1 as ID to avoid conflicts
				fmt.Sprintf("Title %d", index),
				fmt.Sprintf("Content %d", index),
				fmt.Sprintf("Author %d", index),
			)
			require.NoError(t, err)

			err = repo.Create(post)
			assert.NoError(t, err)

			// Try to read the post
			retrievedPost, err := repo.GetByID(post.ID)
			assert.NoError(t, err)
			assert.Equal(t, post.ID, retrievedPost.ID)
		}(i)
	}

	wg.Wait()

	// Verify all posts were created
	posts, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, posts, numGoroutines)
}
