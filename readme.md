# Movie App - Go

1. Command to run the app
    - `go run cmd/api/*.go`

2. Create binary `https://go.dev/doc/install/source#environment`
    - `env GOOS=linux GOARCH=amd64 go build -o gomovies ./cmd/api` 
    - `env GOOS=darwin GOARCH=amd64 go build -o gomovies ./cmd/api` 