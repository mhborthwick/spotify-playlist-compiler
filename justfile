# run auth server
auth:
    @echo 'running on http://localhost:1337'
    go run auth/cmd/main.go

gen:
    @echo 'generating config...'
    pkl-gen-go pkl/Config.pkl --base-path github.com/mhborthwick/spotify-playlist-compiler

run:
    @echo 'running...'
    go run cmd/main.go

build:
    @echo 'building...'
    go build -o bin/spotify-playlist-compiler cmd/main.go
