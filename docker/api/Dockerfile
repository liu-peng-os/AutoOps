FROM golang:1.24-alpine AS builder

ARG GOPROXY=https://proxy.golang.org,direct
ARG GOSUMDB=sum.golang.org

WORKDIR /src/api

ENV GOPROXY=${GOPROXY}
ENV GOSUMDB=${GOSUMDB}

RUN apk add --no-cache ca-certificates git tzdata build-base

COPY api/go.mod api/go.sum ./

RUN go mod download

COPY api/ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -p 1 -o /out/devops main.go

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates openssh-client tzdata curl bash && \
    mkdir -p /app/log /app/upload /root/.ssh

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=builder /out/devops /app/devops

EXPOSE 8000

CMD ["/app/devops"]
