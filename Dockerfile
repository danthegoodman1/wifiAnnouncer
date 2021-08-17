FROM golang:1.16-alpine

WORKDIR /app

COPY go.* .
COPY creds.json .
COPY config.yml .
RUN go mod download

COPY . .

RUN go build -o /app

CMD [ "/app/wifiAnnouncer" ]
