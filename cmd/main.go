package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/mhborthwick/spotify-playlist-compiler/config"
)

type Track struct {
	URI string `json:"uri"`
}

type GetPlaylistResponse struct {
	Items []struct {
		Track Track `json:"track"`
	} `json:"items"`
}

func GetSpotifyURIs(body []byte) ([]string, error) {
	var parsed GetPlaylistResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}
	uris := make([]string, len(parsed.Items))
	for i, item := range parsed.Items {
		uris[i] = item.Track.URI
	}
	return uris, nil
}

func GetSpotifyId(playlist string) (string, error) {
	re := regexp.MustCompile(`[a-zA-Z0-9]{22}`)
	id := re.FindString(playlist)
	if id == "" {
		return "", errors.New("Invalid playlist")
	}
	return id, nil
}

func main() {
	cfg, err := config.LoadFromPath(context.Background(), "pkl/local/config.pkl")
	if err != nil {
		panic(err)
	}

	playlist := "https://open.spotify.com/playlist/4KMuVswhHsgHusA1hSdZQ5?si=a4b8123f214d470c"

	id, err := GetSpotifyId(playlist)

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

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

	uris, err := GetSpotifyURIs(body)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(uris)
}
