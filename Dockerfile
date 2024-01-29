FROM golang:1.21-bullseye as builder

ENV GOOS="linux"
ENV GOARCH="amd64"
ENV CGO_ENABLED="1"
ENV CC="x86_64-linux-gnu-gcc"
ENV CXX="x86_64-linux-gnu-g++"

RUN set -eux && \
    apt-get update && \
    export DEBIAN_FRONTEND=noninteractive && \
    apt-get install -y --no-install-recommends \
        libsqlite3-dev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /vocabacov_build
COPY . .
RUN go build \
    -v \
    -trimpath \
    -buildmode=pie \
    -mod=readonly \
    -tags 'cgo,osusergo,libsqlite3,linux,static,static_build' \
    -ldflags "-s -w -linkmode=external" \
    -o vocabacov ./cmd/vocabacov

FROM alpine:3.19.1

ENV VOCABACOV_DB_PATH="/db/vocabacov.db"

RUN set -eux && \
	apk add --no-cache \
        libc6-compat \
        gcompat \
        sqlite \
        sqlite-libs \
		ca-certificates \
        curl \
        bash && \
    mkdir -p /db && \
    touch /db/vocabacov.db

COPY --from=builder /vocabacov_build/vocabacov /usr/bin/vocabacov

RUN chmod +x /usr/bin/vocabacov

ENTRYPOINT ["/usr/bin/vocabacov"]