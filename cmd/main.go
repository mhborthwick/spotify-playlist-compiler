package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mhborthwick/spotify-playlist-compiler/.config"
)

func main() {
	cfg, err := config.LoadFromPath(context.Background(), "pkl/local/config.pkl")
	if err != nil {
		panic(err)
	}

	id := "4KMuVswhHsgHusA1hSdZQ5"

	url := "https://api.spotify.com/v1/playlists/" + id + "/tracks"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	token := "Bearer " + cfg.Token
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}
