package spotify

import (
	"encoding/json"
	"errors"
	"regexp"
)

// GetID parses the Spotify ID from a playlist link.
func GetID(playlist string) (string, error) {
	re := regexp.MustCompile(`[a-zA-Z0-9]{22}`)
	id := re.FindString(playlist)
	if id == "" {
		return "", errors.New("invalid playlist")
	}
	return id, nil
}

// GetURIs parses the Spotify URIs from a list of tracks.
func GetURIs(body []byte) ([]string, error) {
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
