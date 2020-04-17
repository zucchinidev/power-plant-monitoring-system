FROM golang:1.13.5-stretch as builder

WORKDIR /app/github.com/zucchinidev/power-plant-monitoring-system
COPY . ./
RUN go get ./... && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/sensors sensors/cmd/sensors/*.go
RUN go get ./... && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/coordinators sensors/cmd/coordinators/*.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app/github.com/zucchinidev/power-plant-monitoring-system
COPY --from=builder /go/bin/sensors /usr/bin
COPY --from=builder /go/bin/coordinators /usr/bin
CMD coordinators && sensors

