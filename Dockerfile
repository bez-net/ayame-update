FROM golang:1.13.1-alpine

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ayame

EXPOSE 3000 3443
ENTRYPOINT ["/app/ayame"]