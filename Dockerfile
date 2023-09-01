FROM golang:1.21 AS build

WORKDIR /go/src/futurice/jalapeno

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o bin/jalapeno ./cmd/jalapeno

FROM alpine:3 as production
LABEL org.opencontainers.image.authors="Jalapeno contributors <github.com/futurice/jalapeno>"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL org.opencontainers.image.vendor="Futurice"
LABEL org.opencontainers.image.title="Jalapeno"
LABEL org.opencontainers.image.description="Jalapeno is a CLI for creating, managing and sharing spiced up project templates."
LABEL org.opencontainers.image.url="https://github.com/futurice/jalapeno"
LABEL org.opencontainers.image.source="https://github.com/futurice/jalapeno"
LABEL org.opencontainers.image.documentation="https://futurice.github.io/jalapeno/"

COPY --from=build /go/src/futurice/jalapeno/bin/jalapeno /usr/bin/jalapeno

WORKDIR /workdir

RUN set -eux; \
  addgroup -g 1000 jalapeno; \
  adduser -u 1000 -G jalapeno -s /bin/sh -h /home/jalapeno -D jalapeno

RUN chown -R jalapeno:jalapeno /workdir

USER jalapeno

ENTRYPOINT ["/usr/bin/jalapeno"]
CMD []
