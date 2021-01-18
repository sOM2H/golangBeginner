package main

import (
	"database/sql"
	"encoding/json"
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

func ErrorCheck(err error) {
	if err != nil {
		panic(err.Error())
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
	ErrorCheck(err)
	_, e := insForm.Exec(post.Id, post.UserId, post.Title, post.Body)
	ErrorCheck(e)

	defer db.Close()
}

func insertComment(comment *Comment) {
	db := dbConn()
	insForm, err := db.Prepare("INSERT INTO comments(id, postId, name, email, body) VALUES(?,?,?,?,?)")
	ErrorCheck(err)
	_, e := insForm.Exec(comment.Id, comment.PostId, comment.Name, comment.Email, comment.Body)
	ErrorCheck(e)

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
		go insertPost(&post)
		response2, err := http.Get(commentsURL + strconv.Itoa(post.Id))
		handleError(&err)

		defer response2.Body.Close()

		var comments []Comment
		err = json.Unmarshal(body, &comments)
		for _, comment := range comments {
			go insertComment(&comment)
		}
	}

}
