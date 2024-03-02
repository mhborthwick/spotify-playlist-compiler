gen:
    @echo 'generating config...'
    pkl-gen-go pkl/Config.pkl --base-path github.com/mhborthwick/spotify-playlist-compiler

run:
    @echo 'running...'
    go run cmd/main.go

auth:
    @echo 'running...'
    go run auth/cmd/main.go
