FROM golang:alpine as builder

WORKDIR /app

RUN go install github.com/githubnemo/CompileDaemon@latest

COPY . .

ENTRYPOINT CompileDaemon --build="go build -o auth github.com/sentrionic/ecommerce/auth" --command=./auth