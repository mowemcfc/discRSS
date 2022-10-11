all:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main .
	zip main.zip main .env

arm: 
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o main .
	zip main.zip main .env

clean:
	rm main.zip || true