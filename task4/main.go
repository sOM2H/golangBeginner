package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

func getPost(id int) {
	response, err := http.Get(url + "/" + strconv.Itoa(id))
	handleError(&err)

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	handleError(&err)
	fmt.Println(string(body))
}

func main() {
	for i := 1; i <= 100; i++ {
		go getPost(i)
	}
	var input string
	fmt.Scanln(&input)
}
