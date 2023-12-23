# Go module Download
FROM golang:1.20.3 AS builder
# Go Servide Build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main

# GO service Exec
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs
CMD ["./main"]
