FROM golang:1.19 AS builder

WORKDIR /app
COPY src/ .
RUN go mod download
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app .

FROM alpine:latest

WORKDIR /root/
RUN mkdir /root/logs
COPY --from=builder /app/app ./
EXPOSE 8888
ENTRYPOINT ["./app"]