FROM alpine:3.12

RUN mkdir /app
COPY ./twitchdata /app/twitchdata

ENTRYPOINT ["/app/twitchdata"]