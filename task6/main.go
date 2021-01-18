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

type Post struct {
	Id     int    `json:"id"`
	UserId int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Comment struct {
	Id     int    `json:"id"`
	PostId int    `json:"postId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
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

func insertPost(post *Post, db *sql.DB) {
	insForm, err := db.Prepare("INSERT INTO posts(id, userId, title, body) VALUES(?,?,?,?)")
	if err != nil {
		log.Println(err.Error())
		return
	}
	_, err = insForm.Exec(post.Id, post.UserId, post.Title, post.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}

}

func insertComment(comment *Comment, db *sql.DB) {
	insForm, err := db.Prepare("INSERT INTO comments(id, postId, name, email, body) VALUES(?,?,?,?,?)")
	if err != nil {
		log.Println(err.Error())
		return
	}
	_, err = insForm.Exec(comment.Id, comment.PostId, comment.Name, comment.Email, comment.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}

}

func main() {
	db := dbConn()
	defer db.Close()

	response, err := http.Get(postsURL + strconv.Itoa(7))
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

	var posts []Post
	err = json.Unmarshal(body, &posts)

	for _, post := range posts {
		go insertPost(&post, db)

		response2, err := http.Get(commentsURL + strconv.Itoa(post.Id))
		if err != nil {
			log.Println(err)
			return
		}

		defer response2.Body.Close()

		body, err := ioutil.ReadAll(response2.Body)
		if err != nil {
			log.Println(err)
			return
		}

		var comments []Comment
		err = json.Unmarshal(body, &comments)
		for _, comment := range comments {
			go insertComment(&comment, db)
		}
	}

}
