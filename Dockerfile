FROM golang AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o ./server ./cmd/server

FROM alpine
RUN apk add libc6-compat
WORKDIR /app
COPY .env .env
COPY --from=builder /app/server ./server
CMD ["./server"]