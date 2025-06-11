FROM golang:latest AS builder

COPY . /app

RUN cd /app && go mod tidy && go build -o kaddy .

# multi-stage build
FROM ubuntu:latest

RUN mkdir /app

COPY --from=builder /app/kaddy /app

EXPOSE 8080

WORKDIR /app

CMD ["/app/kaddy"]