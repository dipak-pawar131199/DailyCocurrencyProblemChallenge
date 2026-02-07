package main

import (
	"context"
	"fmt"
	"time"
)

// num producer

func NumProducer(ctx context.Context, nums chan<- int) {
	defer close(nums)
	i := 1
	for i < 50 {

		select {
		case <-ctx.Done():
			fmt.Println("Producer stop")
			return
		case nums <- i:
		}
		i++
	}

}

func NumProcessor(ctx context.Context, nums <-chan int, result chan<- int) {
	defer close(result)
	for num := range nums {
		select {
		case <-ctx.Done():
			fmt.Println("time out for processing")
			return
		case result <- num * num:
			time.Sleep(50 * time.Millisecond)

		}
	}

}
func NumConsumer(ctx context.Context, result <-chan int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("consumer stopped due to timeout")
			return

		case val, ok := <-result:
			if !ok {
				fmt.Println("consumer finished reading all results")
				return
			}
			fmt.Println("Consumed:", val)
		}
	}
}

func main() {
	ctx, cancle := context.WithTimeout(context.Background(), 700*time.Millisecond)

	defer cancle()
	numbers := make(chan int)
	result := make(chan int)
	// Producer
	go NumProducer(ctx, numbers)

	go NumProcessor(ctx, numbers, result)

	go NumConsumer(ctx, result)

	<-ctx.Done()
	fmt.Println("Pipeline stop cleanly")
}
