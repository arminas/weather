package main

import (
	"fmt"
	"log"
	"net/http"
)

const LISTEN_PORT = 8081

func main() {
	http.HandleFunc("/", homepage)
	http.HandleFunc("/favicon.ico", doNothing)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", LISTEN_PORT), nil))
}

func doNothing(w http.ResponseWriter, r *http.Request) {}
