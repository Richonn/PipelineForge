FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server .


FROM alpine:3.23

WORKDIR /app

COPY --from=builder /app/server /server

EXPOSE 8080

CMD [ "/server" ]
