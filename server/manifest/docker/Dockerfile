FROM golang:1.23-alpine AS builder

WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
RUN go build -o family-order-server main.go

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/family-order-server /app/family-order-server
COPY server/manifest/config/config.example.yaml /app/manifest/config/config.yaml
RUN sed -i 's/tcp(127.0.0.1:5432)/tcp(postgres:5432)/' /app/manifest/config/config.yaml
EXPOSE 8000
CMD ["/app/family-order-server"]
