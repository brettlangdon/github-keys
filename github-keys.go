package main

import (
	"log"
	"net/http"
)

var cached []string = make([]string, 0)
var ttl int = 0

func handle(w http.ResponseWriter, r *http.Request) {

}

func main() {

	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
