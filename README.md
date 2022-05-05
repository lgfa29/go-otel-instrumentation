# go-otel-instrumentation

This is a fork of [dalfonzo-tc/go-otel-instrumentation](https://github.com/dalfonzo-tc/go-otel-instrumentation).

I've included the following modifications:
* Containerized the client and server apps
* Included both a GRPC and HTTP version for connecting to the OTel Collector.

This repo is a companion repo to [avillela/hashiqube](https://github.com/avillela/hashiqube), which includes a Nomad jobspec of the OpenTelemetry Collector.

The server and client components of this application send telemetry data to an instance of the OpenTelemetry Collector running in Nomad on HashiQube.

For the complete tutorial, please see my [Medium article](https://adri-v.medium.com/4eaf009b8382?source=friends_link&sk=a1a0612a156d20e86549bd925d419bc3).

## Docker Instructions

This applies to the client and server connecting to the OTel Collector via HTTP only (`client.go` and `server.go`)

1. Build the server and the client image

    ```bash
    docker build -f server.dockerfile -t otel-example-server:1.0.0 .

    docker build -f client.dockerfile -t otel-example-client:1.0.0 .
    ```

2. Run the container instances

    >**NOTE:** Please allow the server time to start up before starting the client.

    ```bash
    # Run server
    docker run -it -p 9000:9000 \
        -h go-server otel-example-server:1.0.0

    # Run client
    docker run -it --rm \
        --network="host" -h go-client \
        otel-example-client:1.0.0
    ```

## Running from the Command-Line

If you have a Go environment set up on your machine and don't want to bother with running the containerized versions, you can run the client and server programs via the command line.

Make sure that you open a separate terminal window for each.

### Running the HTTP version

This connects to the OTel Collector via HTTP.

```bash
go run server.go
```

and in a separate terminal window

```bash
go run client.go
```

### Running the gRPC version

This connects to the OTel Collector via gRPC. Note that I'm using Traefik for load-balancing, and I'm mapping the default OTel Collector gRPC port `4318` to port `7233` on Traefik.

```bash
go run server-grpc.go
```

and in a separate terminal window

```bash
go run client-grpc.go
```
