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

type Error struct {
	Errcode string `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")

	if err != nil {
		log.Printf(err.Error())
	}
	defer db.Close()

	var article Article
	err = db.QueryRow("SELECT id, title, content, created_at, updated_at FROM articles WHERE id = ?", id).Scan(&article.ID, &article.Title, &article.Content, &article.CreatedAt, &article.UpdatedAt)
	if err != nil {
		errResponse := Error{Errcode: "ER1001", Errmsg: "No Data Exists"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAllArticles)
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
