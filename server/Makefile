.PHONY: build
build:
	go build -o main -v ./cmd

.PHONY: build-win
build-win:
	go build -o main.exe -v ./cmd

.PHONY: deploy
deploy:
	scp -r * ${VM_CLOUD_INTERSECTION_ROOT_PATH}/server

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build
