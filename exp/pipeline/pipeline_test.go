package pipeline_test

import (
	"context"
	"slices"
	"sync"
	"testing"
	"time"

	. "github.com/cmilhench/x/exp/pipeline"
)

type Result[T any] struct {
	Out T
	Err error
}

func Input(t *testing.T) <-chan int {
	t.Helper()
	out := make(chan int)
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	delay := []int{400, 400, 400, 100, 0, 0, 0, 300, 400, 0}
	go func() {
		defer close(out)
		for i, v := range values {
			time.Sleep(time.Duration(delay[i]) * time.Millisecond)
			out <- v
		}
	}()
	return out
}

func AssertExecutionTime(t *testing.T, start time.Time, expected time.Duration, tolerance time.Duration) {
	t.Helper()
	elapsed := time.Since(start)
	if elapsed < expected-tolerance || elapsed > expected+tolerance {
		t.Errorf("Execution time was %v, expected around %v (±%v)", elapsed, expected, tolerance)
	}
}

func AssertResults(t *testing.T, result, expected []int) {
	t.Helper()
	if len(result) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(result))
	}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("at index %d, expected %d, got %d", i, v, result[i])
		}
	}
}

func TestGenerator(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	input := []int{1, 2, 3, 4, 5}
	output := Generator(ctx, input...)

	var result []int
	for val := range output {
		result = append(result, val)
	}

	if len(result) != len(input) {
		t.Fatalf("expected length %d, got %d", len(input), len(result))
	}

	for i, v := range input {
		if result[i] != v {
			t.Errorf("at index %d, expected %d, got %d", i, v, result[i])
		}
	}
}

func TestDrop(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	input := Input(t)

	output := Drop(ctx, input, 4)
	var result []int
	for val := range output {
		result = append(result, val)
	}

	AssertExecutionTime(t, start, 2000*time.Millisecond, 10*time.Millisecond)
	AssertResults(t, result, []int{5, 6, 7, 8, 9, 10})
}

func TestTake(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	input := Input(t)

	output := Take(ctx, input, 5)
	var result []int
	for val := range output {
		result = append(result, val)
	}

	AssertExecutionTime(t, start, 1300*time.Millisecond, 10*time.Millisecond)
	AssertResults(t, result, []int{1, 2, 3, 4, 5})
}

func TestFilter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	input := Input(t)

	output := Filter(ctx, input, func(val int) bool {
		return val%2 != 0
	})
	var result []int
	for val := range output {
		result = append(result, val)
	}

	AssertExecutionTime(t, start, 2000*time.Millisecond, 10*time.Millisecond)
	AssertResults(t, result, []int{1, 3, 5, 7, 9})
}

func TestMap(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	input := Input(t)

	output := Map(ctx, 0, input, func(ctx context.Context, id int, val int) int {
		return val * 2
	})
	var result []int
	for val := range output {
		result = append(result, val)
	}

	AssertExecutionTime(t, start, 2000*time.Millisecond, 10*time.Millisecond)
	AssertResults(t, result, []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20})
}

func TestFanIn(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	start := time.Now()
	input1 := make(chan int)
	input2 := make(chan int)
	go func() {
		defer close(input1)
		defer close(input2)
		for i := 0; i < 5; i++ {
			input1 <- i
			input2 <- i * 2
			time.Sleep(200 * time.Millisecond)
		}
	}()

	output := FanIn(ctx, input1, input2)
	var result []int
	for val := range output {
		result = append(result, val)
	}
	slices.Sort(result)

	AssertExecutionTime(t, start, 1000*time.Millisecond, 10*time.Millisecond)
	AssertResults(t, result, []int{0, 0, 1, 2, 2, 3, 4, 4, 6, 8})
}

func TestFanOut(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	input := Input(t)

	workerCount := 3
	var mu sync.Mutex
	processed := make(map[int]int)
	fn := func(ctx context.Context, v int) int {
		mu.Lock()
		processed[v]++
		mu.Unlock()
		time.Sleep(250 * time.Millisecond)
		return v * 2
	}

	outputs := FanOut(ctx, input, workerCount, fn)
	results := make([][]int, workerCount)
	for i := range outputs {
		results[i] = make([]int, 0)
	}
	wg := sync.WaitGroup{}
	for i, ch := range outputs {
		wg.Add(1)
		go func(index int, c <-chan int) {
			defer wg.Done()
			for val := range c {
				results[index] = append(results[index], val)
			}
		}(i, ch)
	}
	wg.Wait()
	var result []int
	for _, v := range results {
		result = append(result, v...)
	}
	slices.Sort(result)

	AssertExecutionTime(t, start, 2500*time.Millisecond, 10*time.Millisecond)
	AssertResults(t, result, []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20})

	for i := 1; i <= 10; i++ {
		if processed[i] != 1 {
			t.Errorf("expected value %d to be processed once, but got %d times", i, processed[i])
		}
	}
}

func TestThrottle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	input := Input(t)

	output := Throttle(ctx, input, 100*time.Millisecond)
	var result []int
	for val := range output {
		result = append(result, val)
	}

	AssertExecutionTime(t, start, 2300*time.Millisecond, 10*time.Millisecond)
	AssertResults(t, result, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
}

func TestPartition(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	input := Input(t)

	outputA, outputB := Partition(ctx, 0, input, func(ctx context.Context, id int, val int) bool {
		return val%2 == 0
	})

	wg := sync.WaitGroup{}
	var resultsA []int
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range outputA {
			resultsA = append(resultsA, v)
		}
	}()

	var resultsB []int
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range outputB {
			resultsB = append(resultsB, v)
		}
	}()

	wg.Wait()

	AssertExecutionTime(t, start, 2000*time.Millisecond, 10*time.Millisecond)
	AssertResults(t, resultsA, []int{2, 4, 6, 8, 10})
	AssertResults(t, resultsB, []int{1, 3, 5, 7, 9})
}

func TestSample(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	input := Input(t)

	output := Sample(ctx, input, 300*time.Millisecond)
	var result []int
	for val := range output {
		result = append(result, val)
	}

	AssertExecutionTime(t, start, 2000*time.Millisecond, 10*time.Millisecond)
	AssertResults(t, result, []int{1, 2, 3, 8, 9})
}

func TestBroadcast(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	input := Input(t)

	numOutputs := 3
	outputs := Broadcast(ctx, input, numOutputs)
	results := make([][]int, numOutputs)
	for i := range outputs {
		results[i] = make([]int, 0)
	}
	wg := sync.WaitGroup{}
	for i, ch := range outputs {
		wg.Add(1)
		go func(index int, c <-chan int) {
			defer wg.Done()
			for val := range c {
				results[index] = append(results[index], val)
			}
		}(i, ch)
	}
	wg.Wait()

	AssertExecutionTime(t, start, 2000*time.Millisecond, 20*time.Millisecond)
	for _, result := range results {
		AssertResults(t, result, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	}
}

func TestChunk(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	input := Input(t)

	output := Chunk(ctx, input, 3)
	var result [][]int
	for vals := range output {
		result = append(result, vals)
	}

	AssertExecutionTime(t, start, 2000*time.Millisecond, 10*time.Millisecond)
	AssertResults(t, result[0], []int{1, 2, 3})
	AssertResults(t, result[1], []int{4, 5, 6})
	AssertResults(t, result[2], []int{7, 8, 9})
	AssertResults(t, result[3], []int{10})
}

func TestWindow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	input := Input(t)

	output := Window(ctx, input, 3, 2)
	var result [][]int
	for vals := range output {
		result = append(result, vals)
	}

	AssertExecutionTime(t, start, 2000*time.Millisecond, 10*time.Millisecond)
	AssertResults(t, result[0], []int{1, 2, 3})
	AssertResults(t, result[1], []int{3, 4, 5})
	AssertResults(t, result[2], []int{5, 6, 7})
	AssertResults(t, result[3], []int{7, 8, 9})
	AssertResults(t, result[4], []int{9, 10})
}

func TestDistinct(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	input := Generator(ctx, 1, 2, 3, 4, 5, 5, 6, 7, 7, 7, 1, 2, 8, 8, 9, 10, 10)

	output := Distinct(ctx, input, func(val int) int { return val })
	var result []int
	for val := range output {
		result = append(result, val)
	}

	AssertResults(t, result, []int{1, 2, 3, 4, 5, 6, 7, 1, 2, 8, 9, 10})
}

func TestRateLimiter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	start := time.Now()
	input := Input(t)

	output := RateLimiter(ctx, input, 3, time.Second, 2)
	var result []int
	for val := range output {
		result = append(result, val)
	}

	AssertExecutionTime(t, start, 3000*time.Millisecond, 10*time.Millisecond)
	AssertResults(t, result, []int{1, 2, 3, 4, 5, 6, 7, 9, 10})
}

func TestReduce(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	input := Input(t)

	output := Reduce(ctx, input, func(acc int, val int) int {
		return acc + val
	})
	var result []int
	for val := range output {
		result = append(result, val)
	}

	AssertExecutionTime(t, start, 2000*time.Millisecond, 10*time.Millisecond)
	AssertResults(t, result, []int{55})
}
