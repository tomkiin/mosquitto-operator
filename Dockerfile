FROM golang:1.18-alpine as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOPROXY=https://goproxy.cn GOARCH=amd64 go build -o manager main.go

FROM alpine:3.15.0
WORKDIR /app
COPY --from=builder /app/manager .

ENTRYPOINT ["/app/manager"]