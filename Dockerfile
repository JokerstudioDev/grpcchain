FROM golang:1.9-alpine AS builder
RUN apk --no-cache add git

RUN go get -d \
    golang.org/x/net/context \
    google.golang.org/grpc \
    github.com/golang/protobuf/proto

COPY . /go/src/pilot.go.grpc
WORKDIR /go/src/pilot.go.grpc/server
RUN go build

FROM alpine:latest
EXPOSE 8080/tcp

WORKDIR /root/
COPY --from=builder /go/src/pilot.go.grpc/server/server .

ENTRYPOINT ["./server"]