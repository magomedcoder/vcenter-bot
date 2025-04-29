.PHONY: run
run:
	go run ./cmd -c=./configs/config.yaml

.PHONY: build
build:
	go build -o ./build/vcenter-bot ./cmd