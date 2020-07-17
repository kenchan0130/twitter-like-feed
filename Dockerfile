FROM golang:latest as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /app
COPY . .
RUN go build main.go

# 作ったGoのバイナリを実行する
FROM alpine:latest
COPY --from=builder /app /app

CMD /app/main $PORT
