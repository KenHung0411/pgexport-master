FROM golang:1.12.6-alpine3.10 as builder

# install upx dep
RUN apk add --no-cache \
    bash \
    binutils \
    ca-certificates \
    curl \
    git \
    tzdata \
    upx \
 && go get -u github.com/golang/dep/...

# setup the working directory
WORKDIR /go/src/gitlab.com/navyx/tools/pgexport

# install dependencies
ADD Gopkg.toml Gopkg.toml
ADD Gopkg.lock Gopkg.lock
RUN dep ensure -vendor-only

# add source code
ADD . .

# build the source
RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` \
    go build -a -installsuffix cgo \
    -ldflags '-s -w -extldflags "-static"' \
    -o /go/bin/pgexport ./cmd/pgexport
RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` \
    go build -a -installsuffix cgo \
    -ldflags '-s -w -extldflags "-static"' \
    -o /go/bin/migrator ./cmd/migrator

# strip and compress the binary
RUN strip --strip-unneeded /go/bin/pgexport \
 && upx -qq /go/bin/pgexport \
 && upx -qq -t /go/bin/pgexport
RUN strip --strip-unneeded /go/bin/migrator \
 && upx -qq /go/bin/migrator \
 && upx -qq -t /go/bin/migrator

# FROM scratch
FROM alpine:3.10

ENV PATH /app:$PATH

WORKDIR /app
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/pgexport /app/pgexport
COPY --from=builder /go/bin/migrator /app/migrator
COPY ./config/config.yml.docker /app/config.yml
COPY ./config/database.yml.docker /app/database.yml
COPY ./migrations /app/migrations

# VOLUME /app/conf

# add launch shell command
COPY docker-entrypoint.sh /usr/bin/

ARG BUILDNO=local
ARG REV=unknown
RUN echo $BUILDNO > $APP_HOME/BUILD; echo $REV > $APP_HOME/REV && \
  echo "INFO: BLUEX-RELEASE DB Exporter --  Rev: $REV -- Build: $BUILDNO" > /MANIFEST

EXPOSE 3000
EXPOSE 5000

ENTRYPOINT ["docker-entrypoint.sh"]
