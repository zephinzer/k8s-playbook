FROM golang:1.12.5-alpine3.9 AS build
RUN apk update --no-cache
RUN apk upgrade --no-cache
RUN apk add --no-cache git
WORKDIR /go/bin
COPY ./bin /go/bin
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN chmod +x /go/bin/pinger

FROM scratch
COPY --from=build /go/bin/pinger /pinger
ENTRYPOINT ["/pinger"]
