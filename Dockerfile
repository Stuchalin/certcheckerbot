FROM golang:1.18-alpine

WORKDIR /app
RUN apk add build-base
RUN apk add sqlite

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY certinfo/*.go ./certinfo/
COPY botprocessing/*.go ./botprocessing/
COPY storage/*.go ./storage/
COPY storage/sqlite3/*.go ./storage/sqlite3/
COPY scheduler/*.go ./scheduler/
COPY *.go ./

RUN go build -o /certcheckerbot
RUN go test ./...

RUN apk del build-base

CMD [ "/certcheckerbot" ]