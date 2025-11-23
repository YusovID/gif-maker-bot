FROM debian:stable-slim

RUN apt-get update && apt-get install -y \
    libstdc++6 \
    libssl3 \
    zlib1g \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY ./telegram-bot-server .

ENTRYPOINT ["./telegram-bot-server"]