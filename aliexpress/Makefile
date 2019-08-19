.PHONY: build
build:
	GOVERSION=`go version` goreleaser --rm-dist

.PHONY: build-snap
build-snap:
	GOVERSION=`go version` goreleaser --snapshot --rm-dist -v
