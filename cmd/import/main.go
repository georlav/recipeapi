package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"sync"

	"github.com/georlav/recipeapi/internal/handler"
)

func main() {
	data, err := ioutil.ReadFile("recipes.json")
	if err != nil {
		log.Fatal(err)
	}

	var recipeRequests []handler.RecipeCreateRequest
	if err := json.Unmarshal(data, &recipeRequests); err != nil {
		log.Fatal(err)
	}

	// Create a channel of recipe request data
	recipesCH := func() chan handler.RecipeCreateRequest {
		ch := make(chan handler.RecipeCreateRequest)

		go func() {
			for i := range recipeRequests {
				ch <- recipeRequests[i]
			}
			close(ch)
		}()

		return ch
	}()

	wgSize := 4
	wg := sync.WaitGroup{}
	wg.Add(wgSize)
	imported := 0

	// Start 4 goroutines
	for i := 1; i <= wgSize; i++ {
		go func() {
			defer wg.Done()

			for r := range recipesCH {
				payload, err := json.Marshal(r)
				if err != nil {
					log.Println(err)
					continue
				}

				resp, err := http.Post(
					"http://127.0.0.1:8080/api/recipes", "application/json", bytes.NewReader(payload),
				)
				if err != nil {
					log.Println(err)
					continue
				}

				if resp.StatusCode != http.StatusCreated {
					respErr, err := httputil.DumpResponse(resp, true)
					if err != nil {
						log.Fatal(err)
					}
					log.Println(string(respErr))
					continue
				}
				resp.Body.Close()

				log.Println("Imported ", r.Title)
				imported++
			}
		}()
	}

	wg.Wait()

	fmt.Printf("File has %d recipes\n", len(recipeRequests))
	fmt.Printf("Imported %d recipes\n", imported)
}
