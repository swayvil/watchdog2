# From https://github.com/chemidy/smallest-secured-golang-docker-image/blob/master/Dockerfile
############################
# STEP 1 build executable binary
############################
# golang alpine 1.14.1
FROM golang:alpine AS builder

# Install git + SSL ca certificates
# Git is required for fetching the dependencies
# Ca-certificates is required to call HTTPS endpoints
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
#ENV USER=appuser
#ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735
#RUN adduser \
#    --disabled-password \
#    --gecos "" \
#    --home "/nonexistent" \
#    --shell "/sbin/nologin" \
#    --no-create-home \
#    --uid "${UID}" \
#    "${USER}"
WORKDIR $GOPATH/src/watchdog2/mail-crawler/
COPY . .

# Fetch dependencies
RUN go get -d -v

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
      -o /go/bin/mail-crawler .

# Create the mounting point
RUN mkdir /store

############################
# STEP 2 build App image
############################
FROM scratch

# Import from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /store /store
#COPY --from=builder /etc/passwd /etc/passwd
#COPY --from=builder /etc/group /etc/group

# Copy our static executable
COPY --from=builder /go/bin/mail-crawler /go/bin/mail-crawler
COPY ./config.json /go/bin/config.json

# Use an unprivileged user
#USER appuser:appuser

# https://stackoverflow.com/questions/59094236/error-unknown-time-zone-america-los-angeles-in-time-loadlocation
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip

EXPOSE 9999

# Run the binary
WORKDIR /go/bin
ENTRYPOINT ["/go/bin/mail-crawler"]