FROM golang:1.20 AS build

ARG taskfileVersion="3.28.0"
ARG taskfileChecksum="dfec009264d35411f893bc0618e924f82bb5188679f88de75169a2d20a1c34f5"

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
LABEL org.opencontainers.image.description="Jalapeno is a CLI for creating, managing and sharing spiced up project templates."
LABEL org.opencontainers.image.url="https://github.com/futurice/jalapeno"
LABEL org.opencontainers.image.source="https://github.com/futurice/jalapeno"
LABEL org.opencontainers.image.documentation="https://futurice.github.io/jalapeno/"

COPY --from=build /work/bin/jalapeno /usr/local/bin/

ENTRYPOINT [ "/usr/local/bin/jalapeno" ]
CMD []