# go-otel-instrumentation

## Open Telemetry: Hands-on Instrumentation

### Instrument once, integrate a thousand times

This repository holds the source code for used in the [OpenTelemetry: Hand-on Example](https://medium.com/@dalfonzo/opentelemetry-hands-on-instrumentation)

## How to use it

1. clone the repository
2. start the otel collector with the default configuration

    ```text
    docker run --rm -p 1111:1111 \
    -v "${PWD}/instrumentation-collector-config.yml":/etc/otel/otel-local-config.yaml \
    --name otelcol otel/opentelemetry-collector-contrib:0.40.0 \
    --config /etc/otel/otel-local-config.yaml \
    "/otelcontribcol"
    ```

3. Start the go server

    ```text
    go run server.go
    ```

4. Send the request using the client

    ```text
    go run client.go 
    ```

5. View the traces in the otel collector output
6. modify the configuration to your liking and restart from #1.
7. Have fun!
