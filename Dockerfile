FROM golang:1.26 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/adora-coding-assessment ./cmd/adora-coding-assessment

FROM alpine:latest
WORKDIR /app

RUN mkdir -p /app/data

COPY --from=builder /app/bin/adora-coding-assessment .

EXPOSE 8080

CMD ["./adora-coding-assessment"]
