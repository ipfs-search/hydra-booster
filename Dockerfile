FROM golang:1.14.2-alpine AS build

RUN apk add --no-cache openssl-dev git build-base

WORKDIR /hydra-booster

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod go.sum ./

#This is the ‘magic’ step that will download all the dependencies that are specified in
# the go.mod and go.sum file.
# Because of how the layer caching system works in Docker, the  go mod download
# command will _ only_ be re-run when the go.mod or go.sum file change
# (or when we add another docker instruction this line)
RUN go mod download && go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get -v

# Copy the source from the current directory
# to the Working Directory inside the container
COPY . .

# Run the build and install
RUN go install -tags=openssl -v ./...

# Create single-layer run image
FROM alpine
COPY --from=build /go/bin/hydra-booster /hydra-booster

RUN apk add --no-cache openssl

# HTTP API
EXPOSE 7779

# Prometheus /metrics
EXPOSE 8888

# Heads
EXPOSE 30000-32767
EXPOSE 30000-32767/udp

CMD ["./hydra-booster", "-metrics-addr=0.0.0.0:8888", "-httpapi-addr=0.0.0.0:7779", "-ui-theme=none"]
