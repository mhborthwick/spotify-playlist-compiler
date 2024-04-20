package spotify

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPlaylistItems(t *testing.T) {
	t.Run("returns body and nil", func(t *testing.T) {
		mockResponse := []byte(``)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(mockResponse)
		}))
		defer mockServer.Close()
		cfg := &Config{
			Token: "YOUR_TOKEN",
		}
		httpClient := &http.Client{}
		data, err := GetPlaylistItems(cfg, httpClient, "123", mockServer.URL)
		assert.Equal(t, mockResponse, data)
		assert.Nil(t, err)
	})
}

func TestCreatePlaylist(t *testing.T) {
	t.Run("returns id and nil", func(t *testing.T) {
		mockResponse := []byte(`{"id": "123"}`)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(mockResponse)
		}))
		defer mockServer.Close()
		cfg := &Config{
			Token: "YOUR_TOKEN",
		}
		httpClient := &http.Client{}
		data, err := CreatePlaylist(cfg, httpClient, mockServer.URL)
		assert.Equal(t, "123", data)
		assert.Nil(t, err)
	})
}

func TestAddItemsToPlaylist(t *testing.T) {
	t.Run("returns body and nil", func(t *testing.T) {
		mockResponse := []byte(``)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(mockResponse)
		}))
		defer mockServer.Close()
		cfg := &Config{
			Token: "YOUR_TOKEN",
		}
		httpClient := &http.Client{}
		data, err := AddItemsToPlaylist(cfg, httpClient, "123", []string{"abc", "def"}, mockServer.URL)
		assert.Equal(t, mockResponse, data)
		assert.Nil(t, err)
	})
}
