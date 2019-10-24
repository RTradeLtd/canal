
GO111MODULE=on

build: fmt
	go build -o canal/canal ./canal

fmt:
	gofmt -w *.go canal/main.go