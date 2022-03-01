FROM golang AS build
RUN apt-get update && apt-get upgrade -y
WORKDIR /go/src/github.com/emvi/oogway
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-s -w" /go/src/github.com/emvi/oogway/cmd/main.go

FROM alpine
RUN apk update && \
    apk upgrade && \
    apk add --no-cache ca-certificates && \
    rm -rf /var/cache/apk/*

# TODO install sass
# https://github.com/sass/dart-sass-embedded/releases/download/1.49.9/sass_embedded-1.49.9-linux-x64.tar.gz

COPY --from=build /go/src/github.com/emvi/oogway/main /oogway/server
RUN addgroup -S oogwayuser && \
    adduser -S -G oogwayuser oogwayuser && \
    chown -R oogwayuser:oogwayuser /oogway
USER oogwayuser
WORKDIR /oogway
VOLUME ["/oogway/data"]
EXPOSE 8080
ENTRYPOINT ["/oogway/server", "/oogway/data"]
