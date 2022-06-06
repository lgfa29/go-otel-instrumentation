# Builder image.
FROM golang:1.17 AS build

ENV CGO_ENABLED=0

WORKDIR /workdir

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o otel-instrumentation .

# Final image.
FROM alpine:3

COPY --from=build /workdir/otel-instrumentation /

EXPOSE 9000
ENTRYPOINT ["/otel-instrumentation"]
