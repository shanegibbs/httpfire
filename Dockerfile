FROM golang:1.17.7-alpine3.15 as builder

WORKDIR /work

ADD go.* .
ADD cmd cmd
ADD pkg pkg
RUN go build -o httpfire ./cmd/agent/main.go

FROM alpine:3.15
COPY --from=builder /work/httpfire /
CMD /httpfire
