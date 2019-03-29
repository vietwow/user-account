FROM golang:alpine

#Install dev tools
RUN apk add --update --no-cache alpine-sdk bash ca-certificates \
      libressl \
      tar \
      git

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine:latest
RUN apk --no-cache add ca-certificates librdkafka

WORKDIR /root/

COPY --from=0 /app/user-account .

EXPOSE 8080
CMD ["./user-account"]