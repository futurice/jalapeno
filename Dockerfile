FROM golang:1.19 AS build

ARG taskfileVersion="3.20.0"
ARG taskfileChecksum="75cd08890ff18f6036255c7630aa17f9ea81bd6b6166747a3913bbb1cff1357c"

WORKDIR /tmp
RUN wget https://github.com/go-task/task/releases/download/v${taskfileVersion}/task_linux_amd64.deb \
	&& echo "${taskfileChecksum}  task_linux_amd64.deb" > checksums.txt \
	&& sha256sum --check checksums.txt \
	&& dpkg -i task_linux_amd64.deb \
	&& rm -f task_linux_amd64.deb

WORKDIR /work
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . ./
RUN task build


FROM debian:stable-slim
LABEL org.opencontainers.image.authors="Jalapeno contributors <github.com/futurice/jalapeno>"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL org.opencontainers.image.vendor="Futurice"
LABEL org.opencontainers.image.title="Jalapeno"
LABEL org.opencontainers.image.url="https://github.com/futurice/jalapeno"
LABEL org.opencontainers.image.source="https://github.com/futurice/jalapeno"
LABEL org.opencontainers.image.documentation="https://futurice.github.io/jalapeno/"

COPY --from=build /work/bin/jalapeno /usr/local/bin/