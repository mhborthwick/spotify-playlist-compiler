gen:
    @echo 'generating config...'
    pkl eval -f yaml config.pkl

run:
    @echo 'running...'
    go run cmd/main.go
