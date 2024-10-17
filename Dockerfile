FROM golang:alpine AS build

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build .

FROM busybox
COPY --from=build /app/dns2ip-templater /sbin/dns2ip-templater

WORKDIR /conf

ENTRYPOINT ["/sbin/dns2ip-templater"]
