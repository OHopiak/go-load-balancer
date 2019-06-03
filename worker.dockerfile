FROM golang:1.11 as builder
WORKDIR $GOPATH/src/github.com/OHopiak/fractal-load-balancer
ENV GO111MODULE on
ENV CGO_ENABLED 0
COPY cmd cmd
COPY core core
COPY worker worker
COPY go.mod go.mod
RUN go install ./cmd/worker

FROM alpine:latest
RUN apk --no-cache add ca-certificates
ENV WORKER_IP "0.0.0.0"
ENV WORKER_PORT "80"
ENV MASTER_IP "master"
ENV MASTER_PORT "8000"
ENV DB_DRIVER "postgres"
ENV DB_CONN_STRING ""
WORKDIR /root/
COPY --from=builder /go/bin/worker .
EXPOSE 80
CMD ["sh", "-c", "sleep 5 && ./worker"]
