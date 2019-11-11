
GO111MODULE=on

VERSION=0.0.22
USER_GH=eyedeekay

version:
	gothub release -s $(GITHUB_TOKEN) -u $(USER_GH) -r canal -t v$(VERSION) -d "Privacy-Enhanced VPN"

build: fmt
	go build -o canal/canal ./canal

fmt:
	gofmt -w *.go canal/main.go