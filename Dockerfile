FROM golang:1.19.0

WORKDIR /usr/src/app

RUN go install github.com/cosmtrek/air@v1.49.0


COPY . .
RUN go mod tidy