package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
)

var numTests, concurrency int
var server string

func main() {
	flag.IntVar(&numTests, "n", 100000, "Number of tests")
	flag.IntVar(&concurrency, "c", 100, "Number of concurrent tests")
	flag.StringVar(&server, "server", "http://localhost:8000", "Server address")
	flag.Parse()

	serverDef := CatServerDefinition{baseUrl: server, resourcesDir: "test_resources"}
	Soak(serverDef, numTests, concurrency)
}

func Soak(def EndpointDefinition, numTests int, concurrency int) {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
	}

	var wg sync.WaitGroup
	wg.Add(concurrency)
	queue := make(chan TestScenario, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for scenario := range queue {
				scenario.Run(client)
			}
		}()
	}
	for i := 0; i < numTests; i++ {
		queue <- def.NextTest()
	}
	close(queue)
	wg.Wait()
}

type EndpointDefinition interface {
	NextTest() TestScenario
}
type TestScenario interface {
	Run(client *http.Client)
}

type CatServerDefinition struct {
	baseUrl      string
	resourcesDir string
}

func (csd CatServerDefinition) NextTest() TestScenario {
	return CreateProfileScenario{CatServerDefinition: csd, profileName: "MissKitty", pictureName: "hai"}
}

type CreateProfileScenario struct {
	CatServerDefinition
	profileName string
	pictureName string
}

func (cps CreateProfileScenario) Run(client *http.Client) {
	url := fmt.Sprintf("%s/%s/", cps.baseUrl, cps.profileName)
	req, err := http.NewRequest("POST", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("POST ", url, "Expected response status 200, got ", resp.StatusCode)
		return
	}

	file, err := os.Open(path.Join(cps.resourcesDir, "cat.jpg"))
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	url = fmt.Sprintf("%s/%s/%s", cps.baseUrl, cps.profileName, cps.pictureName)
	req, err = http.NewRequest("POST", url, nil)
	resp, err = client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("POST ", url, "Expected response status 200, got ", resp.StatusCode)
		return
	}
}
