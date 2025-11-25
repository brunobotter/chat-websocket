# Etapa 1: build
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o server ./cmd/server 

# Etapa 2: runtime
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/server .
COPY cmd/server/dev.yaml /app/dev.yaml 

EXPOSE 8080
CMD ["./server", "--port=8080"]
