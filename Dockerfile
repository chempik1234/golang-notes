# Этап сборки
FROM golang:1.23.3 AS builder
WORKDIR /go/src/notes_service
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/simpleRest ./cmd/simpleRest

# Этап выполнения
FROM debian:buster-slim
WORKDIR /app
COPY --from=builder /go/src/notes_service/bin/simpleRest /app/simpleRest
RUN chmod +x /app/
EXPOSE 3000/tcp
CMD ["/app/simpleRest"]
