PKGS := $(shell go  list ./... | grep -v /vendor)

.PHONY: test
test:
	go test -cover $(PKGS)

angmar:
	CGO_ENABLED=0 go build -o bin/angmar pkg/main/main.go

.PHONY: angmar_stripped
angmar_stripped:
	go build -o bin/angmar -ldflags="-s -w" pkg/main/main.go

.PHONY: angmar_compressed
angmar_compressed: angmar_stripped
	upx bin/angmar