package spotify

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

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

// GetPlaylistItems gets the items (tracks) within a Spotify playlist.
func GetPlaylistItems(cfg *Config, client *http.Client, id string) ([]byte, error) {
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

// CreatePlaylist creates a new empty Spotify playlist.
func CreatePlaylist(cfg *Config, client *http.Client) (string, error) {
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

// AddItemsToPlaylist adds items (tracks) to a playlist.
func AddItemsToPlaylist(
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
