FROM golang:1.19.4 as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /app
COPY . .
RUN go build main.go

FROM alpine:3.18.0
COPY --from=builder /app /app

ENV GIN_MODE release

CMD /app/main $PORT
