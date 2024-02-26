package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", homepage)
	http.HandleFunc("/favicon.ico", doNothing)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func doNothing(w http.ResponseWriter, r *http.Request) {}
