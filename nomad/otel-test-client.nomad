job "otel-test-client" {
  datacenters = ["dc1"]
  type        = "batch"

  parameterized {
    meta_optional = ["name"]
  }

  group "client" {
    task "client" {
      driver = "docker"

      config {
        image   = "laoqui/otel-test:v1"
        command = "client"
        args    = ["-name=${NOMAD_META_name}"]
      }

      template {
        data        = <<EOF
SERVER_URL=http://{{with nomadService "otel-test-server"}}{{with index . 0}}{{.Address}}:{{.Port}}{{end}}{{end}}
OTEL_COLLECTOR_ADDR={{with nomadService "grpc.otel-collector"}}{{with index . 0}}{{.Address}}:{{.Port}}{{end}}{{end}}
OTEL_COLLECTOR_PROTO=grpc
        EOF
        destination = "local/env"
        env         = true
      }
    }
  }
}
