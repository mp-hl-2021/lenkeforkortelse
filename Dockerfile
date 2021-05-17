FROM golang:1.16.2-alpine3.13 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o lenkeforkortelse-server cmd/lenkeforkortelse-server/main.go

FROM alpine:3.13
COPY --from=builder /build/lenkeforkortelse-server .

ENTRYPOINT [ "./lenkeforkortelse-server" ]
