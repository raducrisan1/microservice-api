FROM alpine:3.8
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY ./microservice-api /app/
WORKDIR /app
EXPOSE 3030
CMD ["./microservice-api"]