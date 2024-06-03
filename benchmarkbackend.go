package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	// curl localhost:3000/api/users/show -d '{"username":"pichu"}' -H 'Content-Type: Application/json'
	req, _ := http.NewRequest("POST", "http://localhost:3001/api/users/show", strings.NewReader(`{"username":"pichua"}`))
	req.Header.Set("Content-Type", "Application/json")
	startTime := time.Now()
	for i := 0; i < 1000; i++ {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		_ = body
		// fmt.Println(string(body))
	}
	fmt.Println("Time taken:", time.Since(startTime))
}
