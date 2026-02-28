.PHONY: build clean install test run

build:
	mkdir -p bin
	go build -o bin/tomato ./cmd/tomato

clean:
	rm -rf bin/

install: build
	cp bin/tomato /usr/local/bin/

test:
	go test ./...

run: build
	./bin/tomato
