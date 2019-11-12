# --- build stage
FROM golang:1.13.3-alpine as builder
ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY *.go ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ayame


## --- final stage as base with ffmpeg
FROM jrottenberg/ffmpeg:4.2-alpine
MAINTAINER Stoney Kang <sikang@teamgrit.kr>

# check ffmpeg built
RUN ffmpeg -buildconf

# install mediainfo and check it
RUN apk add mediainfo
RUN mediainfo

# install ayame
WORKDIR /

COPY --from=builder /app/ayame .
COPY config.yaml .
COPY certs/ ./certs/
COPY asset/ ./asset/
COPY upload/ ./upload/

EXPOSE 3000 3443
ENTRYPOINT ["/ayame"]