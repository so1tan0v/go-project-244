APP_PATH ="./cmd/gendiff/main.go"

build:
	go build -o gendiff $(APP_PATH)

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...

run-simple-example:
	make build && ./gendiff examples/simple/file1.json examples/simple/file2.json

run-complex-example:
	make build && ./gendiff examples/complex/file1.json examples/complex/file2.json

run-plain-simple-example:
	make build && ./gendiff -f plain examples/simple/file1.json examples/simple/file2.json

run-plain-complex-example:
	make build && ./gendiff -f plain examples/complex/file1.json examples/complex/file2.json

run-json-simple-example:
	make build && ./gendiff -f json examples/simple/file1.json examples/simple/file2.json

run-json-complex-example:
	make build && ./gendiff -f json examples/complex/file1.json examples/complex/file2.json

lint:
	 golangci-lint run