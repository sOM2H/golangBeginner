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

func handleError(err *error) {
	if *err != nil {
		log.Println(*err)
		return
	}
}
func main() {
	response, err := http.Get(url)
	handleError(&err)

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	handleError(&err)

	fmt.Println(string(body))
}
