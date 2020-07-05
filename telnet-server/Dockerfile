# Build stage
FROM golang:1.14-alpine AS build-env
RUN apk --no-cache add build-base git gcc
ADD . /
RUN cd / && go build -o telnet-server

# final stage
FROM alpine:latest
WORKDIR /app
ENV TELNET_PORT 2323
ENV METRIC_PORT 9000
COPY --from=build-env /telnet-server /app/
CMD ["./telnet-server"]


