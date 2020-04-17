FROM golang:1.13.5-stretch as builder

WORKDIR /app/github.com/zucchinidev/power-plant-monitoring-system
COPY . ./
RUN go get ./... && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/sensors-api sensors/cmd/sensors-api/*.go
RUN go get ./... && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/coordinators-api sensors/cmd/coordinators-api/*.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app/github.com/zucchinidev/power-plant-monitoring-system
COPY --from=builder /go/bin/sensors-api /usr/bin
COPY --from=builder /go/bin/coordinators-api /usr/bin
CMD coordinators-api && sensors-api

