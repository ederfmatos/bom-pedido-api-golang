FROM golang:1.22-alpine as builder
WORKDIR /app
COPY . .
RUN GOOS=linux go build -o app .

FROM alpine
COPY --from=builder /app .
RUN ls -la
CMD ["./app"]