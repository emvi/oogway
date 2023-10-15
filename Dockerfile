FROM golang AS build
RUN apt-get update && apt-get upgrade -y
WORKDIR /go/src/emvi
COPY . oogway
RUN cd /go/src/emvi/pkg && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" cmd/main.go && \
    mkdir /app && \
    mv main /app/pkg

FROM alpine
RUN apk update && \
    apk upgrade && \
    apk add --no-cache && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/* \
WORKDIR /app
COPY --from=build /app /app
#RUN wget https://github.com/sass/dart-sass-embedded/releases/download/1.55.0/sass_embedded-1.55.0-linux-x64.tar.gz && \
#    tar -xf sass_embedded-1.55.0-linux-x64.tar.gz && \
#    mv sass_embedded/dart-sass-embedded /app && \
#    chmod +x /app/dart-sass-embedded
#ENV PATH="${PATH}:/app"
VOLUME ["/app/data"]
EXPOSE 8080
ENTRYPOINT ["/app/oogway", "run", "/app/data"]
