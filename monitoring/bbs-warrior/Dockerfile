# Build stage
FROM golang:alpine AS build-env
ADD . /
RUN cd / && go build -o bbs-warrior

# final stage
FROM alpine:latest as final
WORKDIR /app
COPY --from=build-env /bbs-warrior /app/
CMD ["./bbs-warrior"]


