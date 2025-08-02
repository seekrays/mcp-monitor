ARG BASEIMAGE=alpine:3.18.4
ARG GOVERSION=1.23.0
ARG LDFLAGS=""

# Build the manager binary
FROM golang:${GOVERSION} as builder
# Copy in the go src
WORKDIR /go/src/github.com/seekrays/mcp-monitor
COPY . .
ARG LDFLAGS
ARG TARGETOS
ARG TARGETARCH

# Build
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="${LDFLAGS}" -a -o mcp-monitor /go/src/github.com/seekrays/mcp-monitor

# Copy the cmd into a thin image
FROM ${BASEIMAGE}
WORKDIR /root
RUN apk add gcompat
COPY --from=builder /go/src/github.com/seekrays/mcp-monitor/mcp-monitor /usr/local/bin/mcp-monitor
ENTRYPOINT ["/usr/local/bin/mcp-monitor"]