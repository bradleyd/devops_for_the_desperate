# Build stage
FROM golang:1.15-alpine AS build-env
RUN apk --no-cache add build-base git gcc
ADD . /
RUN cd / && go build -o telnet-server

# final stage
FROM alpine:latest as final
WORKDIR /app
ENV TELNET_PORT 2323
ENV METRIC_PORT 9000
COPY --from=build-env /telnet-server /app/

ENTRYPOINT ["./telnet-server"]

