FROM golang:1.22.4

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app

RUN mkdir -p ./images

EXPOSE 8080

CMD ["./app"]
