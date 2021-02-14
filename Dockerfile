FROM golang:1.14.2-alpine3.11

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

WORKDIR /app
COPY go.mod /app
COPY go.sum /app
RUN go mod download

COPY . /app
RUN go build -o /issue-timer /app

ENTRYPOINT ["/issue-timer"]
