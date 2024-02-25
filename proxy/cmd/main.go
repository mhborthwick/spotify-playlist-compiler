package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/mhborthwick/spotify-playlist-compiler/.config"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello!")
}

func main() {
	cfg, err := config.LoadFromPath(context.Background(), "pkl/local/config.pkl")
	if err != nil {
		panic(err)
	}
	fmt.Printf("I'm running on host %s\n", cfg.Id)
	http.HandleFunc("/", Handler)
	log.Fatal(http.ListenAndServe(":1337", nil))
}
