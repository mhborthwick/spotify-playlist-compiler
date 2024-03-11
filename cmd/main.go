package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/alecthomas/kong"
	"github.com/apple/pkl-go/pkl"
	"github.com/mhborthwick/spotify-playlist-compiler/pkg/spotify"
)

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
		var cfg spotify.Config
		if err = evaluator.EvaluateModule(context.Background(), pkl.FileSource(CLI.Create.Path), &cfg); err != nil {
			panic(err)
		}

		// fmt.Printf("Got module: %+v", cfg)

		playlist := "https://open.spotify.com/playlist/4KMuVswhHsgHusA1hSdZQ5?si=a4b8123f214d470c"

		id, err := spotify.GetID(playlist)

		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}

		client := &http.Client{}

		body, err := spotify.GetPlaylistItems(&cfg, client, id)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		uris, err := spotify.GetURIs(body)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// fmt.Println(uris)

		playlistID, err := spotify.CreatePlaylist(&cfg, client)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// fmt.Println(playlistID)

		_, err = spotify.AddItemsToPlaylist(&cfg, client, playlistID, uris)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println("Playlist created!")
	default:
		panic(ctx.Command())
	}
}
