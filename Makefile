build:
	go build -o bin/main main.go

generate:
	go generate ./...

update:
	go get -u ./...
	go mod tidy

run:
	go run .

test:
	go test -v ./...

