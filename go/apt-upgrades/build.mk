all: buildLinux

buildLinux:
	$(info Building apt-upgrades for linux)
	@env GOOS=linux GOARCH=amd64 go build -o /Users/syacko/workspace/sty-holdings/utilities/bin/apt-upgrades main.go
