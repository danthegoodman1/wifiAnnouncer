FROM golang:1.16-alpine

WORKDIR /app

RUN apk add gcc pkgconfig sdl2-dev --no-cache

COPY go.* /app/
COPY creds.json /app/
COPY config.yml /app/
RUN go mod download

COPY . .

RUN go build -o /app

CMD [ "/app/wifiAnnouncer" ]
