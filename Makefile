.PHONY: run
run:
	go run ./cmd -c=./build/config.yaml

.PHONY: build
build:
	go build -o ./build/vcenter-bot ./cmd