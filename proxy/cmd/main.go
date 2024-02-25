package main

import (
	"fmt"
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello!")
}

func main() {
	http.HandleFunc("/", Handler)
	log.Fatal(http.ListenAndServe(":1337", nil))
}
