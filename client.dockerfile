FROM golang:1.17.5-buster

RUN apt-get update && apt-get install gawk -y

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /workdir

COPY go.mod go.sum client.go startup.sh ./

RUN go mod tidy && \
    chmod +x startup.sh

CMD ["./startup.sh", "client"]
