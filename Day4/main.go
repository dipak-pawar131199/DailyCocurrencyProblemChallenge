package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func processTask(ctx context.Context, w int, tasks <-chan int, result chan<- string, ticker *time.Ticker, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {

		select {
		case <-ctx.Done():
			fmt.Println("Processing stop")
			return
		case <-ticker.C:
			t := time.Now()
			time.Sleep(200 * time.Millisecond)
			r := task*task + 10
			result <- fmt.Sprintf("Worker %d processed task %d result -> %d in %s", w, task, r, time.Since(t))

		}

	}

}

func main() {
	tasks := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	workers := 3
	// rateLimit := 2

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	t := make(chan int)
	result := make(chan string)
	defer cancel()
	var wg sync.WaitGroup
	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go processTask(ctx, i, t, result, ticker, &wg)
	}

	go func() {
		for _, task := range tasks {
			t <- task
		}
		close(t)
	}()

	go func() {
		wg.Wait()
		close(result)
	}()

	for r := range result {
		fmt.Println(r)
	}

}
