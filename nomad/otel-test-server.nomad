job "otel-test-server" {
  datacenters = ["dc1"]

  group "server" {
    network {
      port "http" {}
    }

    service {
      name     = "otel-test-server"
      port     = "http"
      provider = "nomad"
    }

    task "server" {
      driver = "docker"

      config {
        image   = "laoqui/otel-test:v1"
        command = "server"
        ports   = ["http"]
      }

      template {
        data        = <<EOF
BIND_ADDR=0.0.0.0:{{env "NOMAD_PORT_http"}}
OTEL_COLLECTOR_ADDR={{with nomadService "grpc.otel-collector"}}{{with index . 0}}{{.Address}}:{{.Port}}{{end}}{{end}}
OTEL_COLLECTOR_PROTO=grpc
        EOF
        destination = "local/env"
        env         = true
      }
    }
  }
}
