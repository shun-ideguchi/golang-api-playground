FROM golang:1.22.4-bullseye

ENV GO111MODULE=on

WORKDIR /app

COPY /app /app

RUN go install github.com/air-verse/air@v1.52.3
RUN go install github.com/go-delve/delve/cmd/dlv@latest

CMD ["air","-c",".air.toml"]
