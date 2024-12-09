FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY embed.go ./
COPY migrations/ ./migrations/

RUN go build -o gostream ./cmd/gostream

# Stage live
FROM alpine:latest AS binary
WORKDIR /app
COPY --from=builder /app/gostream .
COPY configs/ ./configs/
ENV PATH="/app:${PATH}"

ENTRYPOINT ["./gostream"]

# server is started by default
CMD ["server",  "serve", "-c", "./configs/gostream.toml"]