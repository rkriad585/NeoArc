FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

COPY neoarc-cli/go.mod neoarc-cli/go.sum ./
RUN go mod download

COPY neoarc-cli/ .

ARG VERSION=v0.0.0
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

RUN CGO_ENABLED=0 go build -ldflags "\
  -s -w \
  -X neoarc/internal/cli.Version=${VERSION} \
  -X neoarc/internal/cli.Commit=${COMMIT} \
  -X neoarc/internal/cli.BuildTime=${BUILD_TIME} \
  -X neoarc/internal/version.Version=${VERSION} \
  -X neoarc/internal/version.Commit=${COMMIT} \
  -X neoarc/internal/version.BuildTime=${BUILD_TIME}" \
  -o /usr/local/bin/neoarc ./cmd/neoarc/

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /usr/local/bin/neoarc /usr/local/bin/neoarc

ENTRYPOINT ["neoarc"]
CMD ["help"]
