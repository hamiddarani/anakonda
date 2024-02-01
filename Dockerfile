FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -v -o server ./cmd/main.go

FROM alpine:3.18 AS production

COPY --from=builder /app/server /app/server

CMD [ "/app/server" ] 