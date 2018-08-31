APPNAME := $(shell basename `pwd`)
VERSION := v$(shell gobump show -r)
SRCS := $(shell find . -name "*.go" -type f )
LDFLAGS := -ldflags="-s -w \
	-extldflags \"-static\""
XBUILD_TARGETS := \
	-os="windows linux darwin" \
	-arch="386 amd64" 
DIST_DIR := dist/$(VERSION)
README := README.md
EXTERNAL_TOOLS := \
	github.com/golang/dep/cmd/dep \
	github.com/mitchellh/gox \
	github.com/tcnksm/ghr \
	github.com/motemen/gobump/cmd/gobump

.PHONY: build
build: $(SRCS)
	go build $(LDFLAGS) -o bin/$(APPNAME) .
	go install

.PHONY: xbuild
xbuild: $(SRCS) bootstrap
	gox $(LDFLAGS) $(XBUILD_TARGETS) --output "$(DIST_DIR)/{{.Dir}}_{{.OS}}_{{.Arch}}/{{.Dir}}"

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
release: archive bootstrap
	ghr $(VERSION) $(DIST_DIR)/

.PHONY: test
test:
	go test -cover ./...

.PHONY: clean
clean:
	-rm -rf bin
	-rm -rf $(DIST_DIR)

# 依存ライブラリの更新
.PHONY: deps
deps: bootstrap
	dep ensure

# 外部ツールのインストール
.PHONY: bootstrap
bootstrap:
	for t in $(EXTERNAL_TOOLS); do \
		echo "Installing $$t ..." ; \
		go get $$t ; \
	done

