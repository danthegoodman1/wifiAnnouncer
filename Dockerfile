FROM golang:1.16-buster as build

WORKDIR /app

RUN apt update
RUN apt install libasound2-dev alsa-utils gcc -y

COPY go.* /app/
COPY creds.json /app/
COPY config.yml /app/

RUN go mod download

COPY . .

RUN go version

RUN go build -o /app/wifiAnnouncer

FROM balenalib/raspberrypi4-64-debian
COPY --from=build /app/wifiAnnouncer /app/
COPY --from=build /app/creds.json /app/
COPY --from=build /app/config.yml /app/

RUN apt update
RUN apt install libasound2-dev alsa-utils gcc -y

CMD [ "/app/wifiAnnouncer" ]
