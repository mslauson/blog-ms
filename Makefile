build:
	go build -o bin/main main.go

update:
	go get -u ./...

run:
	go run .

test:
	go test -v ./...

