job "o11y-platform" {
  datacenters = ["dc1"]
  type        = "service"

  group "jaeger" {
    network {
      port "http" {
        to     = 16686
        static = 16686
      }

      port "grpc" {
        to = 14250
      }
    }

    service {
      name     = "jaeger"
      port     = "http"
      tags     = ["http"]
      provider = "nomad"
    }

    service {
      name     = "jaeger"
      port     = "grpc"
      tags     = ["grpc"]
      provider = "nomad"
    }

    task "jaeger-all-in-one" {
      driver = "docker"

      config {
        image = "jaegertracing/all-in-one:latest"
        ports = ["http", "grpc"]
      }

      resources {
        cpu    = 200
        memory = 100
      }
    }
  }

  group "prometheus" {
    network {
      port "http" {
        to     = 9090
        static = 9090
      }
    }

    service {
      name     = "prometheus"
      port     = "http"
      provider = "nomad"
    }

    task "prometheus" {
      driver = "docker"

      config {
        image   = "prom/prometheus:latest"
        ports   = ["http"]
        volumes = ["local/config/prometheus.yaml:/etc/prometheus/prometheus.yml"]
      }

      resources {
        cpu    = 100
        memory = 64
      }

      template {
        data = <<EOF
scrape_configs:
  - job_name: 'otel-collector'
    scrape_interval: 10s
    static_configs:
      - targets:
        {{ range nomadService "metrics.otel-collector" -}}
        - '{{ .Address }}:{{ .Port }}'
        {{ end }}
EOF

        destination = "local/config/prometheus.yaml"
      }
    }
  }

  group "otel-collector" {
    network {
      port "grpc" {
        to     = 4317
        static = 30080
      }

      port "metrics" {
        to = 8889
      }
    }

    service {
      name     = "otel-collector"
      port     = "grpc"
      tags     = ["grpc"]
      provider = "nomad"
    }

    service {
      name     = "otel-collector"
      port     = "metrics"
      tags     = ["metrics"]
      provider = "nomad"
    }

    task "otel-collector" {
      driver = "docker"

      config {
        image = "otel/opentelemetry-collector:0.51.0"

        entrypoint = [
          "/otelcol",
          "--config=local/config.yaml",
        ]

        ports = ["grpc", "metrics"]
      }

      resources {
        cpu    = 200
        memory = 64
      }

      template {
        data = <<EOF
receivers:
  # Make sure to add the otlp receiver.
  # This will open up the receiver on port 4317
  otlp:
    protocols:
      grpc:
        endpoint: "0.0.0.0:{{env "NOMAD_PORT_grpc"}}"
processors:
extensions:
  health_check: {}
exporters:
  jaeger:
    endpoint: "{{with nomadService "grpc.jaeger"}}{{with index . 0}}{{.Address}}:{{.Port}}{{end}}{{end}}"
    tls:
      insecure: true
  prometheus:
    endpoint: 0.0.0.0:{{env "NOMAD_PORT_metrics"}}
    namespace: "testapp"
  logging:

service:
  extensions: [health_check]
  pipelines:
    traces:
      receivers: [otlp]
      processors: []
      exporters: [jaeger]

    metrics:
      receivers: [otlp]
      processors: []
      exporters: [prometheus, logging]
EOF

        destination = "local/config.yaml"
      }
    }
  }
}
