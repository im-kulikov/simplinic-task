FROM golang:1.11-alpine3.8 as builder

RUN apk add --no-cache git

COPY . /simplinic-task

WORKDIR /simplinic-task

# https://github.com/golang/go/wiki/Modules#how-do-i-use-vendoring-with-modules-is-vendoring-going-away
# go build -mod=vendor
RUN set -x \
    && export REPO=$(git rev-parse --short HEAD) \
    && export VERSION=$(git rev-parse --short HEAD) \
    && export BUILD=$(date -u +%s%N) \
    && export LDFLAGS="-w -s -X ${REPO}/misc.BuildVersion=${VERSION} -X ${REPO}/misc.BuildTime=${BUILD}" \
    && export CGO_ENABLED=0 \
    && go build -mod=vendor -ldflags "${LDFLAGS}" -o /go/bin/serve ./cmd/serve/main.go

# Executable image
FROM scratch

WORKDIR /

COPY --from=builder /go/bin/serve                      /serve
COPY --from=builder /simplinic-task/config.yml         /config.yml
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/serve"]
