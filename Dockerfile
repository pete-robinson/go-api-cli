FROM golang:1.17-alpine
ENV GO111MODULE=on

RUN mkdir -p /build
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN cd /build/cmd/cli && go build -o cli .
WORKDIR /dist

RUN cp /build/cmd/cli/cli .
ENTRYPOINT ["/dist/cli"]