FROM golang:1.19-alpine as builder

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /app

COPY . .

RUN make engine

FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app 


WORKDIR /app

EXPOSE 8080

COPY --from=builder /app /app

CMD /app/main