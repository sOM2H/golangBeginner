package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	url = "http://jsonplaceholder.typicode.com/posts"
)

func main() {
	response, err := http.Get(url)
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

	fmt.Println(string(body))
}
