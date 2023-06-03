.PHONY: all api-amd api-arm clean

all: amd arm

amd: api-amd scanner-amd

arm: api-arm scanner-arm

api-amd:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o api cmd/api/main.go
	zip api.zip api .env

api-arm:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o api cmd/api/main.go
	zip api.zip api .env

scanner-amd:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o scanner cmd/scanner/main.go
	zip scanner.zip scanner .env

scanner-arm: 
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o scanner cmd/scanner/main.go
	zip scanner.zip scanner .env

run-api:
	GOARCH=arm64 go run cmd/api/main.go

run-scanner:
	GOARCH=arm64 go run cmd/scanner/main.go

clean:
	rm scanner.zip api.zip scanner api || true
