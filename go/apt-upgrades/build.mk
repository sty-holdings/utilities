all: buildLinux

buildLinux:
	$(info Building apt-upgrades for linux)
	@env GOOS=linux GOARCH=amd64 go build -o bin/apt-upgrades main.go
