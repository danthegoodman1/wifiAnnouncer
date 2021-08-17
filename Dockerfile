FROM golang:1.16-alpine

WORKDIR /app

COPY go.* /app/
COPY creds.json /app/
COPY config.yml /app/
RUN go mod download

COPY . .

RUN go build -o /app

CMD [ "/app/wifiAnnouncer" ]
