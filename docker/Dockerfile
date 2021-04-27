FROM golang:1.16 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux \
    go build -o /server /app/cmd/blog/blog.go

FROM debian:bullseye-slim

WORKDIR /
COPY --from=builder /server /

EXPOSE 8000
CMD ["./server"]
