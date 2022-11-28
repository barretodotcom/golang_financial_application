FROM golang:1.19.3-alpine3.16

USER root

WORKDIR /go/src/app

RUN ls

COPY . .

EXPOSE 3000

CMD ["go", "run", "main.go"]