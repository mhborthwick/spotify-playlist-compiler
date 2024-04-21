package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/apple/pkl-go/pkl"
	"github.com/mhborthwick/spotify-playlists-compiler/pkg/spotify"
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

		// fmt.Printf("Got module: %+v", cfg)

		client := &http.Client{}

		var all []string

		for _, p := range cfg.Playlists {
			id, err := spotify.GetID(p)

			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			body, err := spotify.GetPlaylistItems(&cfg, client, id, "https://api.spotify.com")

			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			uris, err := spotify.GetURIs(body)

			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			all = append(all, uris...)
		}

		playlistID, err := spotify.CreatePlaylist(&cfg, client, "https://api.spotify.com")

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		_, err = spotify.AddItemsToPlaylist(&cfg, client, playlistID, all, "https://api.spotify.com")

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println("Playlist created in: ", time.Since(startNow))
	default:
		panic(ctx.Command())
	}
}
