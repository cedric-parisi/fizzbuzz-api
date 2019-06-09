FROM golang:1.12.5-alpine3.9 as builder

WORKDIR /go/src/github.com/cedric-parisi/fizzbuzz-api

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./fizzbuzz -a -ldflags '-s' -installsuffix cgo ./cmd/fizzbuzz-api/main.go

# Build the migrations
RUN go build migrations/migrate.go

# Run the application
FROM alpine:3.9

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/
COPY --from=builder /go/src/github.com/cedric-parisi/fizzbuzz-api/fizzbuzz .
COPY --from=builder /go/src/github.com/cedric-parisi/fizzbuzz-api/migrate .
COPY --from=builder /go/src/github.com/cedric-parisi/fizzbuzz-api/swaggerui swaggerui
ADD .env .env
ADD migrations/sql migrations/sql

CMD ["sh", "-c", "./migrate up && ./fizzbuzz"]
