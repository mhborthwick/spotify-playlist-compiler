set dotenv-path := 'auth/.env'

id := env_var('CLIENT_ID')
secret := env_var('CLIENT_SECRET')

# run auth server locally
auth:
    @echo 'running on http://localhost:1337'
    @CLIENT_ID={{id}} CLIENT_SECRET={{secret}} go run auth/cmd/main.go

# run auth server on docker
docker:
    @echo 'running on docker...'
    @echo 'running on http://localhost:1337'
    @docker compose -f auth/compose.yml up -d

# stop auth server on docker
stop:
    @echo 'stopping docker...'
    @docker compose -f auth/compose.yml stop

gen:
    @echo 'generating config...'
    pkl-gen-go pkl/Config.pkl --base-path github.com/mhborthwick/spotify-playlist-compiler

run:
    @echo 'running...'
    go run cmd/main.go

build:
    @echo 'building...'
    go build -o bin/spotify-playlist-compiler cmd/main.go
