build_cli:
	mkdir -p ./bin
	go build -o ./bin/cli ./cmd/cli
	go build -o ./bin/interactive ./cmd/interactive