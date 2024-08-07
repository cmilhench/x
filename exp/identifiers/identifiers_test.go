package identifiers

import (
	"sync"
	"testing"
)

func Test_Creator(t *testing.T) {
	ids := make(chan string)
	store := make(map[string]bool)

	workers := 1024
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(instance int) {
			defer wg.Done()
			gen := Creator(uint16(instance))
			for i := 0; i < 512; i++ {
				ids <- gen()
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(ids)
	}()

	for id := range ids {

		if _, exists := store[id]; exists {
			t.Errorf("Duplicate ID generated: %s", id)
		}
		store[id] = true

		_, _, _, err := Parse(id)
		if err != nil {
			t.Errorf("Failed to parse ID: %s", id)
		}
	}
}
