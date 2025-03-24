# syntax=docker/dockerfile:1
# Create a stage for building the application.
ARG GO_VERSION=1.24
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

# Download dependencies as a separate step to take advantage of Docker's caching.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage bind mounts to go.sum and go.mod to avoid having to copy them into
# the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

# Build the application.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage a bind mount to the current directory to avoid having to copy the
# source code into the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -o /bin/bingo ./cmd/bingo/main.go

################################################################################
FROM alpine:latest AS final

RUN --mount=type=cache,target=/var/cache/apk \
    apk add --no-cache \
    tzdata 

WORKDIR /bingo
COPY --from=build /bin/bingo .
RUN chmod +x ./bingo

CMD [ "/bingo/bingo" ]
