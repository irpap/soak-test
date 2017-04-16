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
	flag.IntVar(&requests, "requests", 1, "Number of requests to send")
	flag.IntVar(&concurrency, "concurrency", 100, "Number of concurrent requests")
	flag.StringVar(&server, "server", "http://localhost:8000", "Server address")
	flag.Parse()

	var wg sync.WaitGroup

	hc := http.Client{}

	sem := make(chan struct{}, concurrency)

	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer func() { <-sem; fmt.Println("Remove") }()
			defer wg.Done()
			testUpload(hc)
		}()
	}
	for i := 0; i < requests; i++ {
		sem <- struct{}{}
		fmt.Println("Add")
	}
	wg.Wait()

}

func testUpload(hc http.Client) {
	fmt.Println("Uploading")
	file, err := os.Open("test_resources/test.txt")
	if err != nil {
		fmt.Println(err)
	}
	req, err := http.NewRequest("POST", server+"/test.txt", file)
	resp, err := hc.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	if resp.StatusCode != 200 {
		fmt.Println("Expected response status 200, got ", resp.StatusCode)
	}
	file.Close()
}
