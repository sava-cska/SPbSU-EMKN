.PHONY: build
build:
	go mod tidy
	swag init -g cmd/main.go
	go build -o main.out -v ./cmd

.PHONY: deploy
deploy:
	rsync -av -e ssh --exclude='*.out' . ${VM_CLOUD_INTERSECTION_ROOT_PATH}/server

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build
