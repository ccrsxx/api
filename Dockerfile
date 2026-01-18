ARG GO_VERSION

FROM golang:${GO_VERSION}-alpine AS build

WORKDIR /app

RUN apk add --no-cache tzdata

COPY go.mod go.sum ./

RUN go mod download

COPY src src

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o main ./src/cmd/api/main.go

FROM scratch as final

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /app/main /main

USER 10001:10001

CMD [ "/main" ]

LABEL org.opencontainers.image.authors="ami@ccrsxx.com" \
    org.opencontainers.image.source="https://github.com/ccrsxx/api-go" \
    org.opencontainers.image.description="My personal API for my projects" \
    org.opencontainers.image.licenses="GPL-3.0"

