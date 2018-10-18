package main

import (
	"log"
	"net/http"
)

func main() {

	images := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", images))

	log.Fatal(http.ListenAndServe(":8081", nil))
}
