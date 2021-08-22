FROM balenalib/raspberrypi4-64-golang

WORKDIR /app

RUN apt update
RUN apt install alsa-base alsa-utils gcc -y

COPY go.* /app/
COPY creds.json /app/
COPY config.yml /app/
RUN go mod download

COPY . .

RUN go version

RUN go build -o /app

CMD [ "/app/wifiAnnouncer" ]
