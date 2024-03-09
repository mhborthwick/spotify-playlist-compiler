FROM golang:1.21-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/main.go ./cmd/main.go
RUN CGO_ENABLED=0 go build -o /bin/spotify-auth-server ./cmd/main.go

FROM alpine:latest

WORKDIR /home/spotify-auth-server
RUN addgroup -S spotify && \
    adduser -S spotify-auth-server -G spotify

USER spotify-auth-server
COPY --from=build /bin/spotify-auth-server .local/bin/spotify-auth-server
ENTRYPOINT [".local/bin/spotify-auth-server"]
