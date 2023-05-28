FROM golang:1.19-alpine AS builder

WORKDIR /work

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod download

RUN go build -o /work/raft main.go

FROM alpine:3.15 AS runner

WORKDIR /app

COPY --from=builder /work/raft /app/raft

COPY ./config.yaml /app/config.yaml

EXPOSE 80

CMD /app/raft