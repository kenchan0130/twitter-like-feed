FROM golang:1.17.7 as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /app
COPY . .
RUN go build main.go

FROM alpine:3.15.3
COPY --from=builder /app /app

ENV GIN_MODE release

CMD /app/main $PORT
