FROM golang:1.11 as builder
ADD backend "/go/src/toussaint/backend"
ADD clients "/go/src/toussaint/clients"
WORKDIR "/go/src/toussaint/clients/telegram"
RUN go build .

FROM centos:7
COPY --from=builder "/go/src/toussaint/clients/telegram/telegram" "/usr/bin/"
COPY "docker/telegram/telegram-entrypoint.sh" "entrypoint.sh"
CMD ["/bin/bash", "entrypoint.sh"]