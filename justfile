gen:
    @echo 'generating config...'
    pkl-gen-go pkl/Config.pkl --base-path github.com/mhborthwick/spotify-playlist-compiler

run:
    @echo 'running...'
    go run cmd/main.go
