package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/alecthomas/kong"
	"github.com/apple/pkl-go/pkl"
)

type Track struct {
	URI string `json:"uri"`
}

type GetPlaylistItemsResponseBody struct {
	Items []struct {
		Track Track `json:"track"`
	} `json:"items"`
}

type CreatePlaylistRequestBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

type CreatePlaylistResponseBody struct {
	ID string `json:"id"`
}

type AddItemsToPlaylistRequestBody struct {
	URIs []string `json:"uris"`
}

func GetSpotifyURIs(body []byte) ([]string, error) {
	var parsed GetPlaylistItemsResponseBody
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

func GetSpotifyPlaylistItems(cfg *Config, client *http.Client, id string) ([]byte, error) {
	url := "https://api.spotify.com/v1/playlists/" + id + "/tracks"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	token := "Bearer " + cfg.Token
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func CreateSpotifyPlaylist(cfg *Config, client *http.Client) (string, error) {
	url := "https://api.spotify.com/v1/users/" + cfg.UserID + "/playlists"
	requestData := CreatePlaylistRequestBody{
		Name:        "New Playlist",
		Description: "Created by Playlist Compiler",
		Public:      false,
	}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	token := "Bearer " + cfg.Token
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var parsed CreatePlaylistResponseBody
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", err
	}
	return parsed.ID, nil
}

func AddItemsToSpotifyPlaylist(
	cfg *Config,
	client *http.Client,
	playlistID string,
	uris []string,
) ([]byte, error) {
	url := "https://api.spotify.com/v1/playlists/" + playlistID + "/tracks"
	requestData := AddItemsToPlaylistRequestBody{
		URIs: uris,
	}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	token := "Bearer " + cfg.Token
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	return body, nil
}

type Config struct {
	UserID    string   `pkl:"userID"`
	Token     string   `pkl:"token"`
	Playlists []string `pkl:"playlists"`
}

var CLI struct {
	Create struct {
		Path string `arg:"" name:"path" help:"Path to pkl file." type:"path"`
	} `cmd:"" help:"Create playlist."`
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "create <path>":
		fmt.Println("Evaluating from: " + CLI.Create.Path)

		evaluator, err := pkl.NewEvaluator(context.Background(), pkl.PreconfiguredOptions)
		if err != nil {
			panic(err)
		}
		defer evaluator.Close()
		var cfg Config
		if err = evaluator.EvaluateModule(context.Background(), pkl.FileSource(CLI.Create.Path), &cfg); err != nil {
			panic(err)
		}

		// fmt.Printf("Got module: %+v", cfg)

		playlist := "https://open.spotify.com/playlist/4KMuVswhHsgHusA1hSdZQ5?si=a4b8123f214d470c"

		id, err := GetSpotifyId(playlist)

		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}

		client := &http.Client{}

		body, err := GetSpotifyPlaylistItems(&cfg, client, id)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		uris, err := GetSpotifyURIs(body)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// fmt.Println(uris)

		playlistID, err := CreateSpotifyPlaylist(&cfg, client)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// fmt.Println(playlistID)

		_, err = AddItemsToSpotifyPlaylist(&cfg, client, playlistID, uris)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println("Playlist created!")
	default:
		panic(ctx.Command())
	}
}
