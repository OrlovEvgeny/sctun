BINARY_STUN := stun
BINARY_CTUN := ctun

.PHONY: linux
linux:
	mkdir -p build/linux
	GOOS=linux GOARCH=amd64 go build -o build/linux/$(BINARY_STUN) ./cmd/stun/main.go
	GOOS=linux GOARCH=amd64 go build -o build/linux/$(BINARY_CTUN) ./cmd/ctun/main.go

.PHONY: darwin
darwin:
	mkdir -p build/osx
	GOOS=darwin GOARCH=amd64 go build -o build/osx/$(BINARY_STUN) ./cmd/stun/main.go
	GOOS=darwin GOARCH=amd64 go build -o build/osx/$(BINARY_CTUN) ./cmd/ctun/main.go

.PHONY: win
win:
	mkdir -p build/win
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o build/win/$(BINARY_STUN).exe ./cmd/ctun/main.go
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o build/win/$(BINARY_CTUN).exe ./cmd/ctun/main.go

.PHONY: build
build:  linux darwin