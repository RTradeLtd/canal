
GO111MODULE=on

VERSION=0.0.25
USER_GH=eyedeekay

version:
	gothub release -s $(GITHUB_TOKEN) -u $(USER_GH) -r canal -t v$(VERSION) -d "Privacy-Enhanced VPN"

build: fmt
	go build -o canal/canal ./canal

fmt:
	gofmt -w *.go simplify/*.go canal/main.go

setup:
	sudo ufw --dry-run reset
	sudo ufw --dry-run default deny incoming
	sudo ufw --dry-run default deny outgoing
	sudo ufw --dry-run allow out on tun0 from any to any
	sudo ufw enable

setup-i2p:
	sudo ufw --dry-run allow out from any to any

unsetup:
	sudo ufw reset
	sudo ufw default deny incoming
	sudo ufw default allow outgoing
	sudo ufw enable
