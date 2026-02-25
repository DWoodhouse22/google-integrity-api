ARG GO_VERSION=1.26.0

FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app

# Copy go module files first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the source
COPY . .

# Build the main package (adjust if needed)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -v -o /run-app ./cmd/server


FROM debian:bookworm-slim

# Install CA certificates (required for HTTPS calls)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /run-app /usr/local/bin/run-app

EXPOSE 8080

CMD ["run-app"]