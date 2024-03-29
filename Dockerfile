FROM golang:1.22-alpine

WORKDIR /app

COPY . /app

RUN go mod download

RUN go build -o SSDC

EXPOSE 8081

ENTRYPOINT [ "./SSDC" ]
