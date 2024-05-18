FROM golang:alpine

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apk update && apk add postgresql

RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN go build -o . ./...

CMD ["./url-shortener"]