FROM golang:1.22-alpine AS builder
ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
RUN apk add --no-progress --no-cache gcc musl-dev
WORKDIR /build
COPY . .
RUN go mod download
RUN go build -tags musl -ldflags '-extldflags "-static"' -o app

FROM scratch
WORKDIR /app
COPY --from=builder /build/app ./app
EXPOSE 8080
ENTRYPOINT ["/app/app"]