package pipeline

// comment the functions in this file

import (
	"context"
	"sync"
	"time"
)

// Generator generates values from a slice and sends them to a channel.
// It stops generating values when the 'done' channel is closed.
func Generator[T any](ctx context.Context, in ...T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for _, v := range in {
			select {
			case out <- v: // out ← %+v
			case <-ctx.Done(): // canceled?
				return
			}
		}
	}()
	return out
}

// Take receives values from an input channel and sends only the first n values
// to the output channel.
func Take[T any](ctx context.Context, in <-chan T, n int) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			select {
			case <-ctx.Done():
				return
			case out <- <-in: // out ← in
			}
		}
	}()
	return out
}

// Drop discards the first n values from the input channel and sends the rest
// to the output channel.
func Drop[T any](ctx context.Context, in <-chan T, n int) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			select {
			case <-ctx.Done():
				return
			case <-in: // in
			}
		}
		for v := range in {
			select {
			case <-ctx.Done():
				return
			case out <- v: // out ← in
			}
		}
	}()
	return out
}

// Filter applies a set of functions to values received from an input channel.
// If all functions return true, the value is sent to the output channel.
// The worker stops when the input channel is closed or the context is canceled.
func Filter[T any](ctx context.Context, in <-chan T, fn ...func(T) bool) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
	outer:
		for v := range in {
			select {
			case <-ctx.Done():
				return
			default:
				for _, fn := range fn {
					if !fn(v) {
						continue outer
					}
				}
				out <- v
			}
		}
	}()
	return out
}

// Map applies a function to values received from an input channel.
// It sends the results of the function to an output channel.
// The function stops when the input channel is closed or the context is canceled.
func Map[I, O any](ctx context.Context, id int, in <-chan I, fn func(context.Context, int, I) O) <-chan O {
	out := make(chan O)
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

// FanOut distributes values from the input channel to multiple worker goroutines.
// Each worker receives and processes values using the provided function `fn`.
func FanOut[T, O any](ctx context.Context, in <-chan T, workerCount int, fn func(context.Context, T) O) []<-chan O {
	outs := make([]chan O, workerCount)
	for i := range outs {
		outs[i] = make(chan O)
	}
	var wg sync.WaitGroup
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-in:
					if !ok {
						return
					}
					o := fn(ctx, v)
					select {
					case outs[workerID] <- o:
					case <-ctx.Done():
						return
					}
				}
			}
		}(i)
	}
	go func() {
		wg.Wait()
		for _, ch := range outs {
			close(ch)
		}
	}()
	// Return the output channels (read-only) to the caller
	result := make([]<-chan O, workerCount)
	for i, ch := range outs {
		result[i] = ch
	}
	return result
}

// Throttle limits the rate at which values are sent from the input channel to the output channel.
// Only one value is sent every `rate` duration.
func Throttle[T any](ctx context.Context, in <-chan T, rate time.Duration) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		ticker := time.NewTicker(rate)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- v:
					<-ticker.C
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out
}

// Partition splits an input channel into two based on a predicate function.
// Values satisfying the predicate go to the first channel, others go to the second.
func Partition[T any](ctx context.Context, id int, in <-chan T, fn func(context.Context, int, T) bool) (<-chan T, <-chan T) {
	a := make(chan T)
	b := make(chan T)
	go func() {
		defer close(a)
		defer close(b)
		for i := range in {
			select {
			case <-ctx.Done():
				return
			default:
				if fn(ctx, id, i) {
					a <- i
				} else {
					b <- i
				}
			}
		}
	}()
	return a, b
}

// Sample sends values from the input channel to the output channel at a specified rate.
// based upon time-based sampling. Unlike throttling (which limits the
// frequency of event processing but may process more than one event in each interval),
// time-based sampling captures exactly one event at the specified interval, regardless
// of how many events pass through during that time.
func Sample[T any](ctx context.Context, in <-chan T, rate time.Duration) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		ticker := time.NewTicker(rate)
		defer ticker.Stop()
		for i := range in {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				out <- i
			default:
			}
		}
	}()
	return out
}

// Broadcast sends each value from the input channel to multiple output channels.
func Broadcast[T any](ctx context.Context, in <-chan T, n int) []<-chan T {
	outs := make([]chan T, n)
	for i := 0; i < n; i++ {
		outs[i] = make(chan T)
	}
	go func() {
		defer func() {
			for _, ch := range outs {
				close(ch)
			}
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				var wg sync.WaitGroup
				wg.Add(n)
				for _, ch := range outs {
					go func(c chan T) {
						defer wg.Done()
						select {
						case c <- val:
						case <-ctx.Done():
						}
					}(ch)
				}

				wg.Wait()
			}
		}
	}()
	// Return the output channels (read-only) to the caller
	result := make([]<-chan T, n)
	for i, ch := range outs {
		result[i] = ch
	}
	return result
}

// Chunk groups values from the input channel into slices of a given size and sends each chunk to the output channel.
// The final chunk may be smaller than the chunk size if there are not enough remaining values.
func Chunk[T any](ctx context.Context, in <-chan T, size int) <-chan []T {
	out := make(chan []T)
	go func() {
		defer close(out)
		var chunk []T
		for v := range in {
			select {
			case <-ctx.Done():
				return
			default:
				chunk = append(chunk, v)
				if len(chunk) == size {
					out <- chunk
					chunk = nil
				}
			}
		}
		if len(chunk) > 0 {
			out <- chunk
		}
	}()
	return out
}

// Window groups incoming values into overlapping windows of a specified size and slide.
func Window[T any](ctx context.Context, in <-chan T, size int, slide int) <-chan []T {
	out := make(chan []T)
	go func() {
		defer close(out)
		var window []T
		for v := range in {
			select {
			case <-ctx.Done():
				return
			default:
				window = append(window, v)
				if len(window) > size {
					window = window[1:]
				}
				if len(window) == size {
					out <- window
					window = window[slide:]
				}
			}
		}
		if len(window) > 0 {
			out <- window
		}
	}()
	return out
}

// Distinct removes consecutive duplicate values from the input channel based on a key function.
func Distinct[T any, K comparable](ctx context.Context, in <-chan T, keyFn func(T) K) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		var last K
		for v := range in {
			select {
			case <-ctx.Done():
				return
			default:
				this := keyFn(v)
				if this != last {
					last = this
					out <- v
				}
			}
		}
	}()
	return out
}

// RateLimiter limits the number of values that can pass through in a given time period buffering up values.
// if the buffer exceeds `size` the value passing though is dropped.
func RateLimiter[T any](ctx context.Context, in <-chan T, limit int, per time.Duration, size int) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		ticker := time.NewTicker(per)
		defer ticker.Stop()
		var buffer []T
		var tokens = limit
	outer:
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				tokens = limit
			case v, ok := <-in:
				if !ok {
					break outer
				}
				if tokens > 0 {
					if len(buffer) > 0 {
						out <- buffer[0]
						buffer = buffer[1:]
						if len(buffer) < size {
							buffer = append(buffer, v)
						}
					} else {
						out <- v
					}
					tokens--
				} else {
					if len(buffer) < size {
						buffer = append(buffer, v)
					}
				}
			}
		}
		// drain the buffer
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				tokens = limit
			default:
				if tokens > 0 {
					if len(buffer) > 0 {
						out <- buffer[0]
						buffer = buffer[1:]
						tokens--
					} else {
						return
					}
				}
			}
		}
	}()
	return out
}

// Reduce applies a reducing function to all values from the input channel and produces a single output.
func Reduce[L, R any](ctx context.Context, in <-chan R, fn func(L, R) L) <-chan L {
	out := make(chan L)
	go func() {
		defer close(out)
		var acc L
		for v := range in {
			select {
			case <-ctx.Done():
				return
			default:
				acc = fn(acc, v)
			}
		}
		out <- acc
	}()
	return out
}
