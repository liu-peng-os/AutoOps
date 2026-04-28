FROM golang:1.24-alpine AS builder

ARG ALPINE_REPO_BASE=https://mirrors.aliyun.com/alpine
ARG GO_PROXY=https://mirrors.aliyun.com/goproxy/,https://goproxy.io,direct
ARG GO_SUMDB=off

RUN printf "%s\n%s\n" \
    "${ALPINE_REPO_BASE}/v3.20/main" \
    "${ALPINE_REPO_BASE}/v3.20/community" > /etc/apk/repositories

WORKDIR /src/api

ENV GOPROXY=${GO_PROXY}
ENV GOSUMDB=${GO_SUMDB}

RUN apk add --no-cache ca-certificates git tzdata build-base

COPY api/ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -p 1 -o /out/devops main.go

FROM alpine:3.20

ARG ALPINE_REPO_BASE=https://mirrors.aliyun.com/alpine

RUN printf "%s\n%s\n" \
    "${ALPINE_REPO_BASE}/v3.20/main" \
    "${ALPINE_REPO_BASE}/v3.20/community" > /etc/apk/repositories

WORKDIR /app

RUN apk add --no-cache ca-certificates openssh-client tzdata curl bash && \
    mkdir -p /app/log /app/upload /root/.ssh

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=builder /out/devops /app/devops

EXPOSE 8000

CMD ["/app/devops"]
