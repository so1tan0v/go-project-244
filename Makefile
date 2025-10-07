APP_PATH ="./cmd/gendiff/main.go"

build:
	go build -o gendiff $(APP_PATH)

run:
	./gendiff

run-go:
	go run $(APP_PATH)

test:
	go test