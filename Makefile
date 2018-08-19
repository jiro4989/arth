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
xbuild: $(SRCS) gox
	gox $(LDFLAGS) --output "$(DIST_DIR)/{{.Dir}}_{{.OS}}_{{.Arch}}/{{.Dir}}"

.PHONY: archive
archive: xbuild
	find $(DIST_DIR)/ -mindepth 1 -maxdepth 1 -a -type d \
		| while read -r d; \
		do \
			cp $(README) $$d/ ; \
		done
	cd $(DIST_DIR) && \
		find . -maxdepth 1 -mindepth 1 -a -type d  \
		| while read -r d; \
		do \
			tar czf $$d.tar.gz $$d; \
		done

.PHONY: release
release: archive ghr
	ghr $(VERSION) $(DIST_DIR)/

.PHONY: test
test:
	go test -cover ./...

.PHONY: clean
clean:
	-rm -rf bin
	-rm -rf $(DIST_DIR)

.PHONY: deps
deps: dep
	dep ensure

# 依存するツール

# パッケージ管理
.PHONY: dep
dep:
ifeq ($(shell which dep 2>/dev/null),)
	go get github.com/golang/dep/cmd/dep
endif

# クロスコンパイル
.PHONY: gox
gox:
ifeq ($(shell which gox 2>/dev/null),)
	go get github.com/mitchellh/gox
endif

# githubにリリース
.PHONY: ghr
ghr:
ifeq ($(shell which ghr 2>/dev/null),)
	go get github.com/tcnksm/ghr
endif
