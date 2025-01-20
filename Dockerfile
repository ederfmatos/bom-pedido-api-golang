FROM golang:1.23-alpine AS builder
RUN apk add --no-progress --no-cache gcc musl-dev ca-certificates && update-ca-certificates
WORKDIR /build
COPY . .
RUN PATH="/go/bin:${PATH}" GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go mod tidy && \
    go mod download && \
    go build -tags musl -ldflags '-s -w -extldflags "-static"' -o app ./cmd/api

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/app ./app
COPY --from=builder /build/.resources ./.resources
EXPOSE 8080
ENTRYPOINT ["/app"]