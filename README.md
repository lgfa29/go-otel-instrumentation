# go-otel-instrumentation

This is a fork of [dalfonzo-tc/go-otel-instrumentation](https://github.com/dalfonzo-tc/go-otel-instrumentation).

I've included the following modifications:
* Containerized the client and server apps
* Changed from GRPC to HTTP

This repo is a companion repo to [avillela/hashiqube](https://github.com/avillela/hashiqube), which includes a Nomad jobspec of the OpenTelemetry Collector.

The server and client components of this application send telemetry data to an instance of the OpenTelemetry Collector running in Nomad on HashiQube.

For the complete tutorial, please see my [Medium article](https://adri-v.medium.com/4eaf009b8382?source=friends_link&sk=a1a0612a156d20e86549bd925d419bc3).

## Instructions

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