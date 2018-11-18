FROM golang:1.11.2-stretch as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go get -v ./...
RUN go build -o microservice-api .
FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/microservice-api /app/
WORKDIR /app

EXPOSE 3030

CMD ["./microservice-api"]