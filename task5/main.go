package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	url = "http://jsonplaceholder.typicode.com/posts"
)

func getPost(id int) {
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
}

func main() {
	for i := 1; i <= 100; i++ {
		go getPost(i)
	}
	var input string
	fmt.Scanln(&input)
}
