FROM golang:1.12.0-alpine AS Builder

WORKDIR /go/src/github.com/MovieStoreGuy/orbit
COPY . .

ENV CGO_ENABLED=0 GOOS=linux GO111MODULE=on
RUN set -x && \
    apk --no-cache add git && \
    go build -ldflags="-s -w" -o /orbit

FROM alpine:3.8

LABEL Author="Sean Marciniak"

COPY --from=Builder /orbit /usr/bin/orbit

RUN apk --no-cache add \
        dumb-init \
        ca-certificates \
        openssh-client && \
    addgroup -S orbit && \
    adduser -S orbit -G orbit


USER orbit
WORKDIR /home/orbit

ENTRYPOINT ["dumb-init", "--", "/usr/bin/orbit"]
CMD ["--help"]