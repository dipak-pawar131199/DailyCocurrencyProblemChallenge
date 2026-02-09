package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type nums struct {
	num int
	id  int
}

type result struct {
	num int
	id  int
}

func numProducer(ctx context.Context, numChan chan<- nums, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 10; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Number producing stop")
			return
		default:
			numChan <- nums{num: i, id: id}
			time.Sleep(time.Millisecond * 100)

		}

	}

}

func Fan_In(out chan result, inputs ...<-chan nums) {
	var wg sync.WaitGroup
	for _, ch := range inputs {
		wg.Add(1)

		go func(ch <-chan nums) {
			defer wg.Done()
			for c := range ch {
				out <- result{num: c.num, id: c.id}
			}

		}(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()

}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chan1 := make(chan nums)
	chan2 := make(chan nums)

	mergeChan := make(chan result)

	var wg sync.WaitGroup

	wg.Add(2)

	go numProducer(ctx, chan1, 1, &wg)
	go numProducer(ctx, chan2, 2, &wg)

	go func() {
		wg.Wait()
		close(chan1)
		close(chan2)

	}()

	Fan_In(mergeChan, chan1, chan2)
	timeout := time.After(2 * time.Second)
	for {

		select {
		case value, ok := <-mergeChan:
			if !ok {
				fmt.Println("All data processed")
				return
			}
			t := time.Now()
			r := value.num*value.num + 5
			time.Sleep(time.Millisecond * 200)
			fmt.Printf("Value : %d Result: %d Producer Id : %d Time: %v\n", value.num, r, value.id, time.Since(t))
		case <-timeout:
			fmt.Println("Timeout reached. Cancelling...")
			cancel()
			return

		}

	}
}
