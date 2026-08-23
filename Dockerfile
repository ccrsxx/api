# syntax=docker/dockerfile:1
# check=skip=InvalidDefaultArgInFrom,CopyIgnoredFile

# ---

ARG GO_VERSION

FROM golang:${GO_VERSION}-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o main ./cmd/api/

# ---

FROM scratch AS final

COPY --from=build /app/main /main

USER 10001:10001

ENTRYPOINT [ "/main" ]

EXPOSE 4000

LABEL org.opencontainers.image.authors="ami@ccrsxx.com" \
    org.opencontainers.image.source="https://github.com/ccrsxx/api" \
    org.opencontainers.image.description="My personal API for my projects" \
    org.opencontainers.image.licenses="GPL-3.0"
