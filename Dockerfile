FROM golang:1.14-alpine
WORKDIR /2ofClubsServer
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main
CMD ["./main"]
