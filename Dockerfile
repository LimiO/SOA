FROM golang:1.18-alpine

WORKDIR /practice/converters/
COPY go.* /practice/converters/

RUN go mod download -x

ADD converters /practice/converters/converters
COPY * ./
COPY converters/* ./

ENTRYPOINT ["go", "run", "main.go"]