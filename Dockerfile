FROM golang:1.18-alpine

WORKDIR /app
RUN apk add build-base

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY certinfo/*.go ./certinfo/
COPY botprocessing/*.go ./botprocessing/
COPY *.go ./

RUN go build -o /certcheckerbot
RUN go test ./...

RUN apk del build-base

CMD [ "/certcheckerbot" ]