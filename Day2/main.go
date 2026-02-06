package main

import (
	"fmt"
	"sync"
	"time"
)

type ResultChan struct {
	workerId int
	job      int
	result   int
}

func jobMulitplication(jobsChan chan int, resultChan chan ResultChan, w int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobsChan {
		resultChan <- ResultChan{workerId: w, job: job, result: job * job}
	}

}

func main() {
	jobs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	numWorkers := 3
	jobsChan := make(chan int, len(jobs))
	resultChan := make(chan ResultChan)

	var wg sync.WaitGroup
	t := time.Now()
	for w := range numWorkers {
		wg.Add(1)
		go jobMulitplication(jobsChan, resultChan, w+1, &wg)
	}

	go func() {
		for job := range jobs {
			jobsChan <- job
		}
		close(jobsChan)
	}()
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	for r := range resultChan {
		fmt.Println("worker", r.workerId, "process job", r.job, "result ->", r.result)
	}
	totalTime := time.Since(t)

	fmt.Println("TotalTime", totalTime)

}
