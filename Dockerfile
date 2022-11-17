FROM golang:latest as builder

WORKDIR /app

COPY ./ ./
RUN go build -o main cmd/main.go

FROM ubuntu:latest as exec

COPY --from=builder /app/main ./
COPY --from=builder /app/resources/config/config.yaml ./

ENTRYPOINT ["/main"]
