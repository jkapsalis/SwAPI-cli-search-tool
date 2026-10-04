BINARY    := swapi
PKG       := ./cmd/swapi
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64
GOFLAGS   := -trimpath -ldflags="-s -w"

.PHONY: build release clean

# Static binary for the current platform.
build:
	CGO_ENABLED=0 go build $(GOFLAGS) -o bin/$(BINARY) $(PKG)

# Static binaries for every platform in PLATFORMS, written to dist/.
release:
	@for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; ext=; \
		[ "$$os" = windows ] && ext=.exe; \
		out=dist/$(BINARY)-$$os-$$arch$$ext; \
		echo "building $$out"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build $(GOFLAGS) -o $$out $(PKG) || exit 1; \
	done

clean:
	rm -rf bin dist
