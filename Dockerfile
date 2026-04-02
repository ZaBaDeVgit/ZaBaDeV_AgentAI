FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o zabadev ./cmd/zabadev

FROM alpine:3.20
COPY --from=builder /app/zabadev /usr/local/bin/zabadev
ENTRYPOINT ["zabadev"]