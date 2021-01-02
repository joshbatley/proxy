FROM golang:1.15.6

WORKDIR /app
COPY ./server ./server
COPY ./migrations ./migrations
COPY ./webapp/build ./build

WORKDIR /app/server

RUN go build -o ../bin/api ./api/main.go
EXPOSE 5000
WORKDIR /app/bin

CMD "./api"
