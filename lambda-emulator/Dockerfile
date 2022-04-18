FROM golang:1.18-alpine3.14 AS dev
RUN apk add --no-cache git
ENV PATH="/usr/local/go/bin:${PATH}"

# Hot-Reloader for development
RUN go install github.com/githubnemo/CompileDaemon@latest

# COPY ./configSchemas.json /shared/

COPY ./ /app
WORKDIR /app

RUN go mod download
ENTRYPOINT /go/bin/CompileDaemon --build="go build -o lambda-emulator" --command="./lambda-emulator"


# Testing container
FROM golang:1.18-alpine3.14 AS test
# required cgo setting to run tests in container
ENV CGO_ENABLED 0 

WORKDIR /app
COPY --from=dev /app .
CMD ["sleep", "1d"]


# Production container
FROM golang:1.18-alpine3.14 AS prod
WORKDIR /app
COPY --from=dev /app/lambda-emulator .
CMD [ "./lambda-emulator" ]