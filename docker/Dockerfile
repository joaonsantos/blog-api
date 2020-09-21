FROM golang:1.15 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app/blogapi /app/cmd/blog/blog.go

FROM alpine:3.12
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/blogapi /usr/bin/blogapi
EXPOSE 8000
CMD ["/usr/bin/blogapi"]
