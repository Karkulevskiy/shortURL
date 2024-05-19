FROM golang:alpine as builder

WORKDIR /build

COPY . .
ENV GOPATH=/
RUN go build -o url-shortener ./cmd/url-shortener/main.go

FROM alpine:latest as production
COPY --from=builder /build/url-shortener /build/url-shortener
ADD wait-for-postgres.sh .

RUN apk update && apk add postgresql
RUN chmod +x wait-for-postgres.sh

CMD ["/build/url-shortener"]