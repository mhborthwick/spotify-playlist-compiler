package main

import (
	"context"
	"fmt"

	"github.com/mhborthwick/spotify-playlist-squasher/appconfig"
)

func main() {
	cfg, err := appconfig.LoadFromPath(context.Background(), "pkl/local/appConfig.pkl")
	if err != nil {
		panic(err)
	}
	fmt.Printf("I'm running on host %s\n", cfg.Host)
}
