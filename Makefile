SRCS := $(shell find . -type f | grep .go)
APPNAME := $(shell basename `pwd`)

.PHONY: build
build: $(SRCS)
	go build -o bin/$(APPNAME) .

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	-rm -f bin/*

.PHONY: setup
setup:
	dep ensure
