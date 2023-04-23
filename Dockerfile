FROM golang:1.18-alpine

WORKDIR /practice/converters/
COPY go.* /practice/converters/

RUN go mod download -x

ADD converters /practice/converters/converters
ADD proxy /practice/converters/proxy
COPY * ./
COPY converters/* ./
ENV GROUP_ADDR 239.0.0.1:12345

ENTRYPOINT ["go", "run", "main.go"]