# syntax=docker/dockerfile:1
FROM golang:1.22 AS build-stage

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o / ./...

# Run the tests in the container
FROM build-stage AS run-test-stage

RUN go test -v ./...

FROM gcr.io/distroless/base-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /docker-compose-webhook /docker-compose-webhook

ENV HTTP_PORT=80

USER nonroot:nonroot

# Run
CMD ["/docker-compose-webhook"]
