FROM golang:1.11 as builder
WORKDIR $GOPATH/src/github.com/OHopiak/fractal-load-balancer
ENV GO111MODULE on
ENV CGO_ENABLED 0
COPY cmd cmd
COPY core core
COPY master master
COPY go.mod go.mod
RUN go install ./cmd/master

FROM alpine:latest
RUN apk --no-cache add ca-certificates
ENV MASTER_IP "0.0.0.0"
ENV MASTER_PORT "80"
ENV DB_DRIVER "postgres"
ENV DB_CONN_STRING ""
WORKDIR /root/
COPY --from=builder /go/bin/master master.app
COPY --from=builder /go/src/github.com/OHopiak/fractal-load-balancer/master/static master/static/
COPY --from=builder /go/src/github.com/OHopiak/fractal-load-balancer/master/templates master/templates/
EXPOSE 80
CMD ["sh", "-c", "sleep 5 && ./master.app"]
