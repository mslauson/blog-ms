build:
	go build -o bin/main main.go

update:
	go get -u ./...
	go mod tidy

generate:
	go generate ./...

run:
	go run .

test:
	go test -v ./...

