# Stage 1 - Binary Build
# BUILD_X args should be passed at build time as docker build args
FROM golang:1.12.7-alpine3.10 AS builder
ARG BIN_VERSION
ARG GIT_COMMIT
ARG GIT_SHA
ARG GIT_TAG
ARG GIT_DIRTY
ENV BIN_OUTDIR=./
ENV BIN_NAME=ports-ahoy
RUN apk update && apk add build-base git libressl-dev
WORKDIR /usr/src/ports-ahoy
# install dependencies in separate docker layer
COPY go.mod .
COPY go.sum .
RUN go mod download
# copy application source and build
COPY ./ .
RUN make static

# Stage 2 - Final Image
# The application should be statically linked
FROM alpine:3.10
RUN apk update \
	&& apk add --no-cache ca-certificates \
	&& rm -rf /var/cache/apk/* \
    && addgroup ports-ahoy \
	&& adduser -D -H -G ports-ahoy ports-ahoy
COPY --from=builder /usr/src/ports-ahoy/ports-ahoy /usr/bin/ports-ahoy
VOLUME ["/etc/ports-ahoy"]
ENTRYPOINT ["ports-ahoy", "server"]
EXPOSE 3000
USER ports-ahoy
