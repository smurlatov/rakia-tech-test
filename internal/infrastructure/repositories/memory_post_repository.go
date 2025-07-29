package repositories

import (
	"rakia-tech-test/internal/domain/entities"
	"rakia-tech-test/internal/domain/repositories"
	"sort"
	"sync"
)

type MemoryPostRepository struct {
	posts  map[int]*entities.Post
	nextID int
	mutex  sync.RWMutex
}

func NewMemoryPostRepository() *MemoryPostRepository {
	return &MemoryPostRepository{
		posts:  make(map[int]*entities.Post),
		nextID: 1,
		mutex:  sync.RWMutex{},
	}
}

func (r *MemoryPostRepository) Create(post *entities.Post) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.posts[post.ID]; exists {
		return repositories.ErrPostExists
	}

	// Create a copy to avoid external modifications
	postCopy := *post
	r.posts[post.ID] = &postCopy

	if post.ID >= r.nextID {
		r.nextID = post.ID + 1
	}

	return nil
}

func (r *MemoryPostRepository) CreatePost(title, content, author string) (*entities.Post, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	id := r.nextID
	r.nextID++

	post, err := entities.NewPost(id, title, content, author)
	if err != nil {
		return nil, err
	}

	postCopy := *post
	r.posts[post.ID] = &postCopy

	resultCopy := *post
	return &resultCopy, nil
}

func (r *MemoryPostRepository) GetByID(id int) (*entities.Post, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	post, exists := r.posts[id]
	if !exists {
		return nil, repositories.ErrPostNotFound
	}

	postCopy := *post
	return &postCopy, nil
}

func (r *MemoryPostRepository) GetAll() ([]*entities.Post, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	posts := make([]*entities.Post, 0, len(r.posts))
	for _, post := range r.posts {
		postCopy := *post
		posts = append(posts, &postCopy)
	}

	// Sort posts by ID in ascending order
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID < posts[j].ID
	})

	return posts, nil
}

func (r *MemoryPostRepository) Update(id int, post *entities.Post) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.posts[id]; !exists {
		return repositories.ErrPostNotFound
	}

	postCopy := *post
	postCopy.ID = id
	r.posts[id] = &postCopy

	return nil
}

func (r *MemoryPostRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.posts[id]; !exists {
		return repositories.ErrPostNotFound
	}

	delete(r.posts, id)
	return nil
}

func (r *MemoryPostRepository) Exists(id int) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, exists := r.posts[id]
	return exists
}

func (r *MemoryPostRepository) LoadData(posts []*entities.Post) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, post := range posts {
		postCopy := *post
		r.posts[post.ID] = &postCopy

		if post.ID >= r.nextID {
			r.nextID = post.ID + 1
		}
	}

	return nil
}
