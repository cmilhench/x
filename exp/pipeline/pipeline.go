package pipeline

// comment the functions in this file

import (
	"context"
	"sync"
)

type Result[T any] struct {
	Out T
	Err error
}

// Generator generates values from a slice and sends them to a channel.
// It stops generating values when the 'done' channel is closed.
func Generator[T any](done <-chan struct{}, in ...T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for _, v := range in {
			select {
			case out <- v: // out ← %+v
			case <-done: // canceled?
				return
			}
		}
	}()
	return out
}

// FanIn combines multiple input channels into a single output channel.
// It reads values from the input channels concurrently and sends them to the output channel.
// The function stops when all input channels are closed and all values have been read.
func FanIn[T any](ctx context.Context, in ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	out := make(chan T)

	// Select from all the channels
	for _, v := range in {
		wg.Add(1)
		go func(c <-chan T) {
			defer wg.Done()
			for v := range c {
				select {
				case out <- v: // out ← %+v
				case <-ctx.Done(): // canceled?
					return
				}
			}
		}(v)
	}

	// Wait for all the reads to complete
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Worker applies a function to values received from an input channel.
// It sends the results of the function to an output channel.
// The worker stops when the input channel is closed or the context is canceled.
func Worker[I, O any](ctx context.Context, id int, in <-chan I, fn func(context.Context, int, I) Result[O]) <-chan Result[O] {
	out := make(chan Result[O])
	go func() {
		defer close(out)
		for i := range in {
			select {
			case <-ctx.Done():
				return
			default:
				out <- fn(ctx, id, i)
			}
		}
	}()
	return out
}
