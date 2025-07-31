FROM golang:1.21-alpine AS builder

RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    upx

WORKDIR /app

COPY --from=deps /go/pkg /go/pkg
COPY --from=deps /go/bin /go/bin

COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a \
    -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -o main \
    ./cmd/main.go

RUN chmod +x main