build:
	go build -o mozeidon-z-messaging .

test:
	go test ./...

all: build
