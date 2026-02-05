package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Job struct {
	Index int
	Url   string
}

type Responce struct {
	Index int
	Value string
}

func checkStatusCodeUrls(urlChan <-chan Job, resChan chan<- Responce, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{Timeout: 2 * time.Second}

	for url := range urlChan {

		req, err := http.NewRequest(http.MethodHead, url.Url, nil)
		if err != nil {
			resChan <- Responce{Index: url.Index, Value: "error"}
			continue
		}

		res, err := client.Do(req)

		if err != nil {
			resChan <- Responce{Index: url.Index, Value: "error"}
			continue

		}
		resChan <- Responce{Index: url.Index, Value: "200"}
		res.Body.Close()

	}

}

func main() {
	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://invalid-url",
		"https://golang.org",
	}
	urlJobs := make(chan Job)
	res := make(chan Responce)
	workers := 2
	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < workers; i++ {
		go checkStatusCodeUrls(urlJobs, res, &wg)

	}

	go func() {
		for index, url := range urls {
			urlJobs <- Job{Index: index, Url: url}
		}
		close(urlJobs)
	}()

	// close results AFTER workers finish
	go func() {
		wg.Wait()
		close(res)
	}()

	result := make([]string, len(urls))
	for value := range res {
		result[value.Index] = value.Value
	}
	fmt.Println(result)
	for index, url := range urls {
		fmt.Println(url, "->", result[index])
	}
}
