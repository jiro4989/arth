APPNAME := $(shell basename `pwd`)
VERSION := v1.0.0
SRCS := $(shell find . -name "*.go" -type f )
LDFLAGS := -ldflags="-s -w \
	-X \"main.Version=$(VERSION)\" \
	-extldflags \"-static\""
DIST_DIR := dist/$(VERSION)
README := README.md

.PHONY: build
build: $(SRCS)
	go build $(LDFLAGS) -o bin/$(APPNAME) .
	go install

.PHONY: xbuild
xbuild: $(SRCS)
	gox $(LDFLAGS) --output "$(DIST_DIR)/{{.Dir}}_{{.OS}}_{{.Arch}}/{{.Dir}}"

.PHONY: archive
archive: xbuild
	find $(DIST_DIR)/ -mindepth 1 -maxdepth 1 -a -type d \
		| while read -r d; \
		do \
			cp $(README) $$d/ ; \
		done
	cd $(DIST_DIR) && \
		find . -type d -maxdepth 1 -mindepth 1 \
		| while read -r d; \
		do \
			tar czf $$d.tar.gz $$d; \
		done

.PHONY: release
release: archive
	ghr $(VERSION) $(DIST_DIR)/

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	-rm -f bin
	-rm -rf $(DIST_DIR)

.PHONY: dep
dep:
ifeq ($(shell which dep 2>/dev/null),)
	go get -u github.com/golang/dep/cmd/dep
endif

.PHONY: deps
deps: dep
	dep ensure
