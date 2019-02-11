FROM golang:1.11 as builder
ADD backend "/go/src/toussaint/backend"
WORKDIR "/go/src/toussaint/backend"
RUN go build .

FROM centos:7
COPY --from=builder "/go/src/toussaint/backend/backend" "/usr/bin"
RUN echo "backend --host=$IP" > /usr/bin/entrypoint.sh
WORKDIR /home
ENTRYPOINT ["/bin/bash", "/usr/bin/entrypoint.sh"]
