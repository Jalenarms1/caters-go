FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .
COPY go.sum go.mod ./
RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./app/caters-go ./cmd/cater-go
RUN apk add --no-cache ca-certificates

FROM scratch
WORKDIR /root

COPY --from=builder ./app/caters-go .

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

CMD ["./caters-go"]