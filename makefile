
build:
	   @go build -o bin/gobank2
run: build
	   @./bin/gobank2

test:
	   @go test -v ./...