APP_PATH ="./cmd/gendiff/main.go"

build:
	go build -o gendiff $(APP_PATH)

run:
	./gendiff $(ARGS)

run-go:
	go run $(APP_PATH)

test:
	go test ./...

run-example:
	make build && ./gendiff examples/file1.json examples/file2.json

lint:
	 golangci-lint run