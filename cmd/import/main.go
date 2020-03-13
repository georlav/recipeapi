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

	"github.com/georlav/recipeapi/internal/database"
)

func main() {
	data, err := ioutil.ReadFile("recipes.json")
	if err != nil {
		log.Fatal(err)
	}

	var recipes database.Recipes
	if err := json.Unmarshal(data, &recipes); err != nil {
		log.Fatal(err)
	}

	// Create a channel with recipes
	recipesCH := func() chan database.Recipe {
		ch := make(chan database.Recipe)

		go func() {
			for i := range recipes {
				ch <- recipes[i]
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
					"http://127.0.0.1:8080/api/v1/recipes", "application/json", bytes.NewReader(payload),
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

	fmt.Printf("File has %d recipes\n", len(recipes))
	fmt.Printf("Imported %d recipes\n", imported)
}
