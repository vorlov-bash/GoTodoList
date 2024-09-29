build_cli:
	mkdir -p ./bin
	go build -o ./bin/cli ./cmd/cli
	go build -o ./bin/http ./cmd/http

run-http:
	make build_cli
	./bin/http