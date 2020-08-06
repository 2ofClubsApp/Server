FROM golang:1.14-alpine
WORKDIR /2ofClubsServer
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main
CMD ["./main"]

#RUN apk update && apk upgrade && \
 #    apk add --no-cache bash git openssh
#FROM postgres
#COPY *.sql /docker-entrypoint-initdb.d/
#ADD init.sql /docker-entrypoint-initdb.d/
#RUN chmod a+r /docker-entrypoint-initdb.d/*
