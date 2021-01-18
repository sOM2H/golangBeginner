package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const (
	url = "http://jsonplaceholder.typicode.com/posts"
)

var wg sync.WaitGroup

func getPost(id int) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		response, err := http.Get(url + "/" + strconv.Itoa(id))
		if err != nil {
			log.Println(err)
			return
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println(err)
			return
		}

		f, err := os.Create("storage/posts/" + strconv.Itoa(id) + ".txt")
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()

		_, err = f.WriteString(string(body))
		if err != nil {
			log.Println(err)
			return
		}
	}()
}

func main() {
	for i := 1; i <= 100; i++ {
		go getPost(i)
	}

	wg.Wait()
}
