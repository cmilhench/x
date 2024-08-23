package identifiers_test

import (
	"sync"
	"testing"

	. "github.com/cmilhench/x/exp/identifiers"
)

func Test_Creator(t *testing.T) {
	ids := make(chan uint64)
	store := make(map[uint64]bool)

	workers := 1024
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := range workers {
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
			t.Errorf("Duplicate ID generated: %d", id)
		}
		store[id] = true

		_, _, _, err := Parse(id)
		if err != nil {
			t.Errorf("Failed to parse ID: %d", id)
		}
	}
}
