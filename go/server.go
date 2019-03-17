package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"github.com/olivere/elastic"
	"bytes"
)

//NumberOfWorkers ...
var NumberOfWorkers = 3

//NumberOfJobs ...
var NumberOfJobs = 20

var active = true
var elasticSearchURL = "http://34.244.29.39:9200"
var serviceURL = "http://0.0.0.0:8080"

//Jobs ...
type Jobs struct {
	data elastic.SearchHit
	id   int
}

//Results ...
type Results struct {
	msg   string
	id    int
	Error bool
}

func main() {
	//Start elasic search client
	client, err := elastic.NewClient(elastic.SetURL(elasticSearchURL), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	
	//Create channel jobs
	jobs := make(chan Jobs, 100)
	//Create channel results
	results := make(chan Results, 100)
	//Start workers
	for w := 1; w <= NumberOfWorkers; w++ {
		go worker(w, jobs, results)
	}

	//fetch elastic search data
	for j:=1; j <= NumberOfJobs; j++ {
		go fetch(jobs, client)
	}
	
	//Get worker results (not active for now)
	for r := range results {
		if r.Error {
			fmt.Println(r.msg)
		} else {
			fmt.Println("jobID: ", r.id)
		}
	}
	close(results)

}

func fetch(ch chan Jobs, client *elastic.Client) {
	scroll := client.Scroll("filebeat*").Size(1000)
	ctx := context.Background()
	for active {
		//get scroll results
		searchResult, err := scroll.Do(ctx)
		
		if err == io.EOF {
			fmt.Println("done")
			active = false
			continue
		}

		if err != nil {
			panic(err)
		}

		if len(searchResult.Hits.Hits) > 0 {
			//push jobs to channel
			for i, job := range searchResult.Hits.Hits {
				ch <- Jobs{*job, i}
			}
		}

	}
	close(ch)
}

func worker(id int, jobs <-chan Jobs, results chan<- Results) {
	fmt.Println("worker ", id, "start running")
	for j := range jobs {
		fmt.Println(j.id)
		sendToNewEndPoint(j)
	}
}

func sendToNewEndPoint(j Jobs) {
	//encode data json
	var jsonStr, err = j.data.Source.MarshalJSON()
	req, err := http.NewRequest("POST", serviceURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return
	}

	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//If other service is not available
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//print service response
	fmt.Println("response Body:", string(body))
}
