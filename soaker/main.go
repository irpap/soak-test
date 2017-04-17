package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
)

var requests, concurrency int
var server string

func main() {
	flag.IntVar(&requests, "n", 100000, "Number of requests to send")
	flag.IntVar(&concurrency, "c", 100, "Number of concurrent requests")
	flag.StringVar(&server, "server", "http://localhost:8000", "Server address")
	flag.Parse()

	var wg sync.WaitGroup

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
	}

	queue := make(chan int, concurrency)
	wg.Add(requests)

	for i := 0; i < concurrency; i++ {
		go func() {
			for _ = range queue {
				// fmt.Println("Executed:", n)
				defer wg.Done()
				testUpload(client)
			}

		}()

	}
	for i := 0; i < requests; i++ {
		queue <- i
		// fmt.Println("Queued: ", i)
	}
	close(queue)
	wg.Wait()
}

func testUpload(client *http.Client) {
	file, err := os.Open("test_resources/test.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	req, err := http.NewRequest("POST", server+"/test.txt", file)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("the request errored: ", err)
	}

	if resp == nil {
		panic("Received nil response")
	}
	if resp.StatusCode != 200 {
		fmt.Println("Expected response status 200, got ", resp.StatusCode)
	}
}
