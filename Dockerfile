FROM golang:1.16 AS builder
WORKDIR /go/src/webdav
COPY . /go/src/webdav
ENV CGO_ENABLED=0
RUN go install

FROM scratch
COPY --from=builder /go/bin/webdav /

ENV LISTEN :8080
ENV ROOT /webdav
ENV PREFIX /
EXPOSE 8080/tcp
ENTRYPOINT ["/webdav"]
