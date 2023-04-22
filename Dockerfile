FROM golang:1.18-alpine

WORKDIR /pracrice/converters/
COPY go.* .

RUN go mod download -x

COPY * .

RUN go build .

ENTRYPOINT ["./main"]