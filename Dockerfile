FROM golang:latest as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /app
COPY . .
RUN go build main.go

FROM alpine:latest
COPY --from=builder /app /app

ENV GIN_MODE release

CMD /app/main $PORT
