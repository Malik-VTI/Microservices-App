# Stage 1 - Build binary
FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o product-service .

# Stage 2 - Run binary
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/product-service .
EXPOSE 8082
ENTRYPOINT ["./product-service"]
