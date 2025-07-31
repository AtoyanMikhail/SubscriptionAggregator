FROM golang:1.24.4-alpine3.21 AS deps

RUN apk add --no-cache \
    git \
    make \
    gcc \
    musl-dev \
    ca-certificates && \
    rm -rf /var/cache/apk/*

FROM deps AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -a \
    -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -o main \
    ./cmd/main.go

RUN chmod +x main

FROM alpine:3.21

# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/main /main

COPY --from=builder /app/configs /configs

COPY --from=builder /app/migrations /migrations

RUN addgroup -S api && \
    adduser -S api -G api && \
    chown api:api /main && \
    mkdir -p /static && \
    chown api:api /static && \
    rm -rf /var/cache/apk/*

USER api

EXPOSE $PORT

ENTRYPOINT ["/main"]