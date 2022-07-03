FROM golang:1.18.2-alpine3.15 AS dev

RUN apk add --no-cache \
    build-base \
    gcc \
    git
    
RUN go install github.com/githubnemo/CompileDaemon@v1.4.0

COPY ./ /app
WORKDIR /app

RUN go mod download
RUN go mod tidy
RUN go build main.go
ENTRYPOINT /go/bin/CompileDaemon --build="go build main.go"



# Production container
FROM golang:1.18-alpine3.14 AS prod
RUN apk add --update docker openrc
RUN rc-update add docker boot
WORKDIR /app
COPY --from=dev /app/main .
CMD [ "./main" ]