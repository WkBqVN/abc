
FROM golang:1.21

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go mod download

COPY . .
COPY .env .

RUN go build -o /main

EXPOSE 8080

CMD ["/main"]