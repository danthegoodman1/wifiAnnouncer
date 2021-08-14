FROM golang:1.16-alpine

WORKDIR /app

COPY go.* .
COPY creds.json .
RUN go mod download

COPY *.go .

RUN go build -o /app

CMD [ "/app/wifiAnnouncer" ]
