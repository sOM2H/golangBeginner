package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

const (
	postsURL    = "http://jsonplaceholder.typicode.com/posts?userId="
	commentsURL = "http://jsonplaceholder.typicode.com/comments?postId="
)

// Post struct
type Post struct {
	Id     int    `json:"id"`
	UserId int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// Comment struct
type Comment struct {
	Id     int    `json:"id"`
	PostId int    `json:"postId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

func handleError(err *error) {
	if *err != nil {
		log.Println(*err)
		return
	}
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "golangbeginner"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func insertPost(post *Post) {
	db := dbConn()
	insForm, err := db.Prepare("INSERT INTO posts(id, userId, title, body) VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(post.Id)
	insForm.Exec(post.Id, post.UserId, post.Title, post.Body)
	defer db.Close()
}

func insertComment(comment *Comment) {
	db := dbConn()
	insForm, err := db.Prepare("INSERT INTO comments(id, postId, name, email, body) VALUES(?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(comment.Id, comment.PostId, comment.Name, comment.Email, comment.Body)
	defer db.Close()

}

func main() {

	response, err := http.Get(postsURL + strconv.Itoa(7))
	handleError(&err)

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	handleError(&err)

	var posts []Post
	err = json.Unmarshal(body, &posts)

	for _, post := range posts {
		insertPost(&post)
		response2, err2 := http.Get(commentsURL + strconv.Itoa(post.Id))
		handleError(&err2)

		defer response2.Body.Close()

		var comments []Comment
		err = json.Unmarshal(body, &comments)
		for _, comment := range comments {
			insertComment(&comment)
		}
	}

}
