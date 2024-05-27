# stage: build ---------------------------------------------------------

FROM golang:1.22-alpine as build

RUN apk add --no-cache gcc musl-dev linux-headers

WORKDIR /go/src/github.com/flashbots/latency-monitor

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o bin/latency-monitor -ldflags "-s -w" github.com/flashbots/latency-monitor/cmd

# stage: run -----------------------------------------------------------

FROM alpine

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=build /go/src/github.com/flashbots/latency-monitor/bin/latency-monitor ./latency-monitor

ENTRYPOINT ["/app/latency-monitor"]
