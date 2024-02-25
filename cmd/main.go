package main

import (
	"context"
	"fmt"

	"github.com/mhborthwick/spotify-playlist-compiler/.config"
)

func main() {
	cfg, err := config.LoadFromPath(context.Background(), "pkl/local/config.pkl")
	if err != nil {
		panic(err)
	}
	fmt.Printf("I'm running on host %s\n", cfg.Host)
}
