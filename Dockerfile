FROM golang:1.17.7-alpine3.15 as builder

WORKDIR /work

ADD go.* .
ADD cmd cmd
ADD pkg pkg
RUN go build -o agent ./cmd/agent/main.go
RUN go build -o director ./cmd/director/main.go

FROM alpine:3.15
COPY --from=builder /work/agent /usr/local/bin/
COPY --from=builder /work/director /usr/local/bin/
