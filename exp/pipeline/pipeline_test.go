package pipeline_test

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"

	. "github.com/cmilhench/x/exp/pipeline"
)

func Test(t *testing.T) {
	// create a timeout mechanism
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	n, err := process(ctx, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 0})
	if err != nil {
		t.Fatal(err)
	}
	if n > runtime.NumCPU() {
		t.Fatalf("expected >%d, actual %d\n", runtime.NumCPU(), n)
	}
}

func process(ctx context.Context, numbers []int) (int, error) {
	workers := runtime.NumCPU()
	// listen for context cancellation (get a channel)
	done := ctx.Done()

	// turn our slice or work in to a channel
	tasks := Generator(done, numbers...)
	// spin up the workers
	results := make([]<-chan Result[int], workers)
	for i := 0; i < workers; i++ {
		results[i] = Map(ctx, i, tasks, divide)
	}
	// receive the answers
	count := 0
	errs := []error{}
	for answer := range FanIn(ctx, results...) {
		count = count + 1

		if answer.Err != nil {
			errs = append(errs, answer.Err)
		}
	}
	return count, errors.Join(errs...)
}

func divide(ctx context.Context, id int, in int) Result[int] {
	time.Sleep(2000 * time.Millisecond)
	return Result[int]{Out: 100 / in}
}

type Result[T any] struct {
	Out T
	Err error
}
