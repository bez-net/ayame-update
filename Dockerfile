# build stage
FROM golang:1.13.3-alpine as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY *.go ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ayame

# final stage
FROM scratch
COPY --from=builder /app/ayame /

COPY config.yaml /
COPY certs/* /certs/
COPY assets/* /assets/

EXPOSE 3000 3443
ENTRYPOINT ["/ayame"]