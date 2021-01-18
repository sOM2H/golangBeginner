package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
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

type Database struct {
	*sql.DB
}

var wg sync.WaitGroup

func (db Database) Cmd(cmd string, a ...interface{}) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		insForm, err := db.Prepare(cmd)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = insForm.Exec(a...)
		if err != nil {
			log.Println(err)
			return
		}
	}()
}

func dbConnect() (Database, error) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "golangbeginner"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		return Database{nil}, err
	}
	return Database{db}, nil
}

func unmarshalResponse(url string, i interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &i)
	return err
}

func main() {
	var posts []Post

	db, err := dbConnect()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	db.Cmd("drop table if exists posts, comments")
	wg.Wait()

	db.Cmd("create table posts(id INT NOT NULL, userId INT NOT NULL, title VARCHAR(100) NOT NULL, body VARCHAR(8000) NOT NULL)")
	db.Cmd("create table comments(id INT NOT NULL, postId INT NOT NULL, name VARCHAR(100) NOT NULL, email VARCHAR(100) NOT NULL, body VARCHAR(8000) NOT NULL)")
	wg.Wait()

	err = unmarshalResponse(postsURL+strconv.Itoa(7), &posts)
	if err != nil {
		log.Println(err)
		return
	}

	for _, post := range posts {
		db.Cmd("INSERT INTO posts(id, userId, title, body) VALUES(?,?,?,?)", post.Id, post.UserId, post.Title, post.Body)

		var comments []Comment
		err = unmarshalResponse(commentsURL+strconv.Itoa(post.Id), &comments)
		if err != nil {
			log.Println(err)
			return
		}

		for _, comment := range comments {
			db.Cmd("INSERT INTO comments(id, postId, name, email, body) VALUES(?,?,?,?,?)", comment.Id, comment.PostId, comment.Name, comment.Email, comment.Body)
		}
	}

	wg.Wait()
}
