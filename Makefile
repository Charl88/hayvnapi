install:
	go clean && \
	go get -d -v ./... && \
	go install -v ./... && \
	go mod tidy
start:
	go run .
build:
	go build .