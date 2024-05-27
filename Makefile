VERSION := $(shell git describe --tags --always --dirty="-dev" --match "v*.*.*" || echo "development" )
VERSION := $(VERSION:v%=%)

.PHONY: build
build:
	@CGO_ENABLED=0 go build \
			-ldflags "-X main.version=${VERSION}" \
			-o ./bin/latency-monitor \
		github.com/flashbots/latency-monitor/cmd

.PHONY: snapshot
snapshot:
	@goreleaser release --snapshot --clean
