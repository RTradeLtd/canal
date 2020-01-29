
GO111MODULE=on

VERSION=0.0.26
USER_GH=eyedeekay

build: fmt
	go build -o canal

version:
	gothub release -s $(GITHUB_TOKEN) -u $(USER_GH) -r canal -t v$(VERSION) -d "VPN Configuration Tool for Go"

fmt:
	gofmt -w *.go etc/*.go

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

readme:
	@echo "canal - VPN Auto-Configuration Tool" | tee README.md
	@echo "===================================" | tee -a README.md
	@echo "" | tee -a README.md
	@echo "This utility encapsulates commands for configuring a VPN as the default route" | tee -a README.md
	@echo "on straightforward client applications, or forwarding to another connection on" | tee -a README.md
	@echo "a server." | tee -a README.md
	@echo "" | tee -a README.md
	./canal -h 2>&1 | tr '\t' ' ' | sed "s|  |          |g" | sed "s|                |        |g" | sed "s|Usage|        Usage|g" | tee -a README.md
	@echo "" | tee -a README.md
	@echo "It is considerably less likely to destroy things you love now, but it probably" | tee -a README.md
	@echo "still needs to be smarter than it is." | tee -a README.md
	@echo "" | tee -a README.md
