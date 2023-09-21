FROM golang AS builder1

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -ldflags "-s -w -extldflags '-static'" -o istio-server

FROM alpine
ENV PORT=8000
COPY --from=builder1 /build/istio-server /
EXPOSE 8000
WORKDIR /data
ENTRYPOINT ["/istio-server"]