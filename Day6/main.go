package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var count int

type errorResult struct {
	err         string
	workerId    int
	num         int
	processTime time.Duration
}

type opResult struct {
	result      int
	workerId    int
	num         int
	processTime time.Duration
}

func numGenerator(ctx context.Context, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 30; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Number Generation stop")
			return
		case out <- i:

		}

	}
}

func processorFAN_OUT(ctx context.Context, wId int, numchan <-chan int, err chan<- errorResult, out chan<- opResult, wg *sync.WaitGroup) {
	defer wg.Done()
	// var mtx sync.Mutex
	for num := range numchan {
		select {
		case <-ctx.Done():
			fmt.Println("Fan Out operation stop...")
			return
		default:
			if num%7 == 0 {
				t := time.Now()
				time.Sleep(time.Millisecond * 20)

				err <- errorResult{err: "Number divisible by 7", workerId: wId, num: num, processTime: time.Since(t)}
				// mtx.Lock()
				// count++
				// mtx.Unlock()

			} else {
				t := time.Now()
				r := num * num
				time.Sleep(time.Millisecond * 20)
				out <- opResult{result: r, workerId: wId, num: num, processTime: time.Since(t)}
			}

		}
	}

}
func main() {
	numWorkers := 4
	maxError := 3
	numchan := make(chan int)
	errchan := make(chan errorResult)
	opResult := make(chan opResult)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg1, wg sync.WaitGroup
	wg.Add(1)
	go numGenerator(ctx, numchan, &wg)

	go func() {
		wg.Wait()
		close(numchan)
	}()

	// for num := range numchan {
	// 	fmt.Println(num)
	// }

	for w := 1; w <= numWorkers; w++ {
		wg1.Add(1)

		go processorFAN_OUT(ctx, w, numchan, errchan, opResult, &wg1)

	}
	go func() {
		wg1.Wait()
		close(errchan)
		close(opResult)
	}()

	errorCounter := 0
	for errchan != nil || opResult != nil {
		select {
		case err, ok := <-errchan:
			if !ok {
				fmt.Println("All err reading from channel done")
				errchan = nil
				continue
			}
			errorCounter++
			fmt.Println("err", err)
			if errorCounter > maxError {
				fmt.Println("Max error reached cancel context")
				// errchan = nil
				cancel()
				// continue
			}
		case out, ok := <-opResult:
			if !ok {
				fmt.Println("All out reading from channel done")
				opResult = nil
				continue
			}
			fmt.Println("out", out)
		}
	}

}
