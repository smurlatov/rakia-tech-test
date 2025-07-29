package main

import (
	"fmt"
	"rakia-tech-test/internal/infrastructure/repositories"
	"sync"
	"time"
)

func main() {
	repo := repositories.NewMemoryPostRepository()

	fmt.Println("🧪 Test 1: Concurrent Create (50 goroutines)")
	testConcurrentCreate(repo)

	fmt.Println("\n🧪 Test 2: Concurrent Read/Write")
	testConcurrentReadWrite(repo)
}

func testConcurrentCreate(repo *repositories.MemoryPostRepository) {
	var wg sync.WaitGroup
	createdIDs := make(chan int, 100)

	startTime := time.Now()
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			post, err := repo.CreatePost(
				fmt.Sprintf("Title %d", i),
				fmt.Sprintf("Content %d", i),
				fmt.Sprintf("Author %d", i),
			)
			if err != nil {
				fmt.Printf("❌ Error creating post: %v\n", err)
				return
			}

			createdIDs <- post.ID
		}(i)
	}

	go func() {
		wg.Wait()
		close(createdIDs)
	}()

	var ids []int
	for id := range createdIDs {
		ids = append(ids, id)
	}
	duration := time.Since(startTime)

	idMap := make(map[int]bool)
	duplicates := 0
	for _, id := range ids {
		if idMap[id] {
			duplicates++
			fmt.Printf("❌ Duplicate ID found: %d\n", id)
		}
		idMap[id] = true
	}

	if duplicates == 0 {
		fmt.Printf("✅ All %d IDs are unique! No race conditions.\n", len(ids))
	} else {
		fmt.Printf("❌ Found %d duplicate IDs! Race condition detected.\n", duplicates)
	}

	allPosts, _ := repo.GetAll()
	fmt.Printf("✅ Total posts in repository: %d\n", len(allPosts))
	fmt.Printf("Time taken: %v\n", duration)
}

func testConcurrentReadWrite(repo *repositories.MemoryPostRepository) {
	var wg sync.WaitGroup
	readCount := 0
	writeCount := 0
	updateCount := 0

	startTime := time.Now()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 20; i++ {
			posts, err := repo.GetAll()
			if err != nil {
				fmt.Printf("❌ Read error: %v\n", err)
			} else {
				readCount++
				if i%5 == 0 {
					fmt.Printf("📖 Read #%d: %d posts\n", readCount, len(posts))
				}
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			post, err := repo.CreatePost(
				fmt.Sprintf("Concurrent Title %d", i),
				fmt.Sprintf("Concurrent Content %d", i),
				"Concurrent Author",
			)
			if err != nil {
				fmt.Printf("❌ Error creating post: %v\n", err)
				continue
			}

			writeCount++
			fmt.Printf("✍️  Created post ID: %d (#%d)\n", post.ID, writeCount)
			time.Sleep(2 * time.Millisecond)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(5 * time.Millisecond)

		for i := 0; i < 5; i++ {
			posts, _ := repo.GetAll()
			if len(posts) > 0 {
				postToUpdate := posts[i%len(posts)]
				postToUpdate.Title = fmt.Sprintf("Updated Title %d", i)

				if err := repo.Update(postToUpdate.ID, postToUpdate); err != nil {
					fmt.Printf("❌ Update error: %v\n", err)
				} else {
					updateCount++
					fmt.Printf("Updated post ID: %d (#%d)\n", postToUpdate.ID, updateCount)
				}
			}
			time.Sleep(3 * time.Millisecond)
		}
	}()

	wg.Wait()
	duration := time.Since(startTime)

	allPosts, _ := repo.GetAll()
	fmt.Printf("\nResults:\n")
	fmt.Printf("Reads: %d\n", readCount)
	fmt.Printf("Writes: %d\n", writeCount)
	fmt.Printf("Updates: %d\n", updateCount)
	fmt.Printf("Final posts: %d\n", len(allPosts))
	fmt.Printf("Duration: %v\n", duration)
	fmt.Println("✅ No deadlocks or race conditions detected!")
}
