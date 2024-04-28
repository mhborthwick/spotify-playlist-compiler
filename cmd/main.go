package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/apple/pkl-go/pkl"
	"github.com/mhborthwick/spotify-playlists-combiner/pkg/spotify"
)

var CLI struct {
	Create struct {
		Path string `arg:"" name:"path" help:"Path to pkl file." type:"path"`
	} `cmd:"" help:"Create playlist."`
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "create <path>":
		startNow := time.Now()
		fmt.Println("Evaluating from: " + CLI.Create.Path)

		evaluator, err := pkl.NewEvaluator(context.Background(), pkl.PreconfiguredOptions)
		if err != nil {
			panic(err)
		}
		defer evaluator.Close()
		var cfg spotify.Config
		if err = evaluator.EvaluateModule(context.Background(), pkl.FileSource(CLI.Create.Path), &cfg); err != nil {
			panic(err)
		}

		spotifyClient := spotify.Spotify{
			URL:    "https://api.spotify.com",
			Token:  cfg.Token,
			UserID: cfg.UserID,
			Client: &http.Client{},
		}

		var all []string

		for _, p := range cfg.Playlists {
			id, err := spotify.GetID(p)
			handleError(err)
			body, err := spotifyClient.GetPlaylistItems(id)
			handleError(err)
			uris, err := spotify.GetURIs(body)
			handleError(err)
			all = append(all, uris...)
		}

		playlistID, err := spotifyClient.CreatePlaylist()
		handleError(err)
		_, err = spotifyClient.AddItemsToPlaylist(all, playlistID)
		handleError(err)

		fmt.Println("Playlist created in: ", time.Since(startNow))
	default:
		panic(ctx.Command())
	}
}
