FROM golang:1-alpine as builder

RUN apk --no-cache --no-progress add git ca-certificates tzdata make \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

WORKDIR /go/acme-fixer

# Download go modules
COPY go.mod .
COPY go.sum .
RUN GO111MODULE=on go mod download

COPY . .

RUN make build

FROM alpine:3.11

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/acme-fixer/acme-fixer .

ENTRYPOINT ["/acme-fixer"]
