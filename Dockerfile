FROM golang AS builder1

ENV GO111MODULE=on \
    GOOS=linux \
    CGO_ENABLED=0 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -ldflags "-s -w -extldflags '-static'" -o istio-server

FROM alpine
ENV PORT=8000
RUN apk update \
    && apk upgrade \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates 2>/dev/null || true
ENV TZ="Asia/Shanghai"
COPY --from=builder1 /build/istio-server /
EXPOSE 8000
ENTRYPOINT ["/istio-server"]