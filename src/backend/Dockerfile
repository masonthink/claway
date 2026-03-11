# Multi-stage build for Claway backend
# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

# Build migrate CLI
RUN CGO_ENABLED=0 go install -tags 'postgres' \
    github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /bin/server \
    ./cmd/server

# Stage 2: Runtime
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /bin/server /server
COPY --from=builder /go/bin/migrate /migrate
COPY --from=builder /app/migrations /migrations
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/server"]
