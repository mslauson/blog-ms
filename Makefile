build:
	go build -o bin/main main.go

update:
	go get -u ./...

generate:
	go generate ./...

run:
	go run .

test:
	go test -v ./...

