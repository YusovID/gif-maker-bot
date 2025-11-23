FROM golang:1.25.4-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/gif-maker

RUN CGO_ENABLED=0 GOOS=linux go build -o /gif-maker .

FROM alpine:3.22

RUN apk update && apk add --no-cache ffmpeg

WORKDIR /

COPY --from=builder /gif-maker /gif-maker
COPY ./config/dev.yml /config/dev.yml

EXPOSE 8080

CMD ["/gif-maker"]