# Building stage
FROM golang:1.13-alpine3.10 as builder
LABEL maintainer="Kyle Bai<k2r2.bai@gmail.com>"

ENV GOPATH "/go"
ENV PROJECT_PATH "$GOPATH/src/github.com/kairen/rotate-elasticsearch-index"
ENV GO111MODULE "on"

RUN apk add --no-cache git ca-certificates g++ make

COPY . $PROJECT_PATH
RUN cd $PROJECT_PATH && \
    make && mv out/rotate-index /tmp/rotate-index

# Running stage
FROM alpine:3.10
COPY --from=builder /tmp/rotate-index /bin/rotate-index
ENTRYPOINT [ "rotate-index" ]