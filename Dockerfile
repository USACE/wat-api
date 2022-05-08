FROM osgeo/gdal:alpine-normal-3.2.1 as dev

COPY --from=golang:1.18-alpine3.14 /usr/local/go/ /usr/local/go/

RUN apk add --no-cache \
    pkgconfig \
    gcc \
    libc-dev \
    git

ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV GO111MODULE="on"
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin
RUN apk add --no-cache git
ENV CGO_ENABLED 0 

# Hot-Reloader for development
RUN go install github.com/githubnemo/CompileDaemon@latest

# COPY ./configSchemas.json /shared/

COPY ./ /app
WORKDIR /app

RUN go mod download
RUN go mod tidy
RUN go build main.go
ENTRYPOINT /go/bin/CompileDaemon --build="go build main.go"


# Testing container
#FROM golang:1.18-alpine3.14 AS test
# required cgo setting to run tests in container
#ENV CGO_ENABLED 0 

#WORKDIR /app
#COPY --from=dev /app .
#CMD ["sleep", "1d"]


# Production container
FROM golang:1.18-alpine3.14 AS prod
RUN apk add --update docker openrc
RUN rc-update add docker boot
WORKDIR /app
COPY --from=dev /app/main .
CMD [ "./main" ]