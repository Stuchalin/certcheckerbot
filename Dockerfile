FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY certinfo/*.go ./certinfo/
COPY *.go ./

RUN go build -o /certcheckerbot

CMD [ "/certcheckerbot" ]