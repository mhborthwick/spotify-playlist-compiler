gen:
    @echo 'generating config...'
    pkl-gen-go pkl/AppConfig.pkl --base-path github.com/mhborthwick/spotify-playlist-squasher

run:
    @echo 'running...'
    go run cmd/main.go
