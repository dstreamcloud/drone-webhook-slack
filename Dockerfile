FROM golang:1.14 AS build

WORKDIR /go/src/github.com/dstream.cloud/drone-webhook-slack
ADD . /go/src/github.com/dstream.cloud/drone-webhook-slack

RUN go build -o /go/bin/drone-webhook-slack main.go

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/drone-webhook-slack /
ENTRYPOINT ["/drone-webhook-slack"]
