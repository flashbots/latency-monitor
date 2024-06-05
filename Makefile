VERSION := $(shell git describe --tags --always --dirty="-dev" --match "v*.*.*" || echo "development" )
VERSION := $(VERSION:v%=%)

.PHONY: build
build:
	@CGO_ENABLED=0 go build \
			-ldflags "-X main.version=${VERSION}" \
			-o ./bin/latency-monitor \
		github.com/flashbots/latency-monitor/cmd

.PHONY: run
run:
	@go run github.com/flashbots/latency-monitor/cmd

.PHONY: local-test
local-test:
	@go run github.com/flashbots/latency-monitor/cmd serve \
		--transponder-peer localhost=127.0.0.1:32123 \
		--transponder-interval 1s

.PHONY: snapshot
snapshot:
	@goreleaser release --snapshot --clean
