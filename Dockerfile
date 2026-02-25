FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /rate-limiter ./cmd/server

FROM alpine:3.19
COPY --from=builder /rate-limiter /rate-limiter
COPY config.yaml /config.yaml
EXPOSE 8090
ENTRYPOINT ["/rate-limiter"]
