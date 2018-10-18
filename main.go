package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Article struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Articles []Article

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Homepage.")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")

	if err != nil {
		log.Printf(err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT id, title, content, created_at, updated_at from articles")
	if err != nil {
		panic(err.Error())
	}

	var articles []Article
	for results.Next() {
		var article Article

		// for each row, scan the result into our article composite object
		err = results.Scan(&article.ID, &article.Title, &article.Content, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}

		articles = append(articles, article)
	}
	json.NewEncoder(w).Encode(articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	articles := Articles{
		Article{ID: 1, Title: "Hello", Content: "Article Content"},
		Article{ID: 2, Title: "Hello 2", Content: "Article Content"},
	}

	for _, v := range articles {
		if v.ID == id {
			json.NewEncoder(w).Encode(v)
		}
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAllArticles)
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	log.Fatal(http.ListenAndServe(":10001", myRouter))
}

func main() {
	handleRequests()
}
