.PHONY: build
build:
	go build -o main -v ./cmd

.PHONY: build-win
build-win:
	go build -o main.exe -v ./cmd

.PHONY: deploy
deploy:
	go build -o main -v ./cmd
	scp -i ${VM_SSH_ID_RSA_PATH} main ${VM_CLOUD_INTERSECTION_ROOT_PATH}
	scp -i ${VM_SSH_ID_RSA_PATH} configs/config.toml ${VM_CLOUD_INTERSECTION_ROOT_PATH}/configs/config.toml

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build
