# 📊 opentelemetry-research

Eksperimen observabilitas menggunakan OpenTelemetry di aplikasi Go, lengkap dengan integrasi Prometheus dan Grafana. Proyek ini mensimulasikan perekaman metrik seperti nilai acak dan penggunaan memori heap secara real-time.

## ✨ Fitur

- Integrasi OpenTelemetry SDK di aplikasi Go
- Ekspor metrik ke Prometheus melalui OpenTelemetry Collector
- Visualisasi metrik menggunakan Grafana
- Error tracing menggunakan Jaeger
- Konfigurasi lengkap menggunakan Docker Compose

## 📦 Struktur Direktori

```
opentelemetry-research/
├── main.go             # Aplikasi Go dengan integrasi OpenTelemetry
├── go.mod / go.sum     # Dependency Go
├── docker-compose.yml  # Layanan Prometheus, Grafana, dan Otel Collector
├── otel-collector-config.yaml # Konfigurasi OpenTelemetry Collector
├── pkg/
│   └── telemetry/
│       └── telemetry.go # Inisiasi dari telemetry
├── prometheus/
│   └── prometheus.yml # Konfigurasi Prometheus
└── grafana/
    └── provisioning/
        └── datasources/
            └── datasource.yml # Auto-provision Prometheus ke Grafana
```

## 🚀 Cara Menjalankan

### 1. Clone Repository

```bash
git clone https://github.com/putuadityabayu/opentelemetry-research.git
cd opentelemetry-research
```

### 2. Jalankan Docker Compose

```bash
docker-compose up -d
```

Docker Compose akan menjalankan:

- Prometheus di port 9090
- Grafana di port 3000
- Jaeger UI di port 16686
- OpenTelemetry Collector (terhubung internal container)

### 3. Jalankan Aplikasi Go

```bash
go run main.go
```

## 📈 Akses dan Visualisasi

- Grafana: http://localhost:3000
  - Login: admin / admin
- Prometheus UI: http://localhost:9090
- Jaeger UI: http://localhost:16686

## 🔧 Konfigurasi Penting

### `otel-collector-config.yaml`

Konfigurasi collector yang menangkap metrik dari aplikasi dan meneruskannya ke Prometheus.

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:
    # Send batch every 10 seconds or when batch size reaches 100 items
    timeout: 10s
    send_batch_size: 100

  memory_limiter:
    check_interval: 1s
    limit_percentage: 80
    spike_limit_percentage: 25

exporters:
  prometheus:
    endpoint: "0.0.0.0:8889"
    namespace: "otel"
    send_timestamps: true
    metric_expiration: 180m
    resource_to_telemetry_conversion:
      enabled: true

  otlp:
    endpoint: "jaeger:4317"
    tls:
      insecure: true

  debug:
    verbosity: detailed

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [otlp, debug]
    metrics:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [prometheus, debug]
```

### `prometheus/prometheus.yml`

Konfigurasi untuk prometheus:

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'otel-collector'
    scrape_interval: 10s
    static_configs:
      - targets: ['otel-collector:8889']
        labels:
          group: 'otel-collector'

  - job_name: 'prometheus'
    scrape_interval: 10s
    static_configs:
      - targets: ['localhost:9090']
        labels:
          group: 'prometheus'
```

### `grafana/provisioning/datasources/datasource.yml`

Konfigurasi otomatis datasource Grafana:

```yaml
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    editable: true

  - name: Jaeger
    type: jaeger
    access: proxy
    url: http://jaeger:16686
    editable: true
```

## 🛠️ Dependencies

Pastikan Anda telah menginstall:
- [Docker](https://www.docker.com/)
- [Go](https://golang.org/) (versi terbaru)

## 📚 Referensi

- [OpenTelemetry Go SDK](https://opentelemetry.io/docs/instrumentation/go/)
- [Prometheus](https://prometheus.io/)
- [Grafana](https://grafana.com/)
- [Jaeger](https://www.jaegertracing.io/)

## 🤝 Kontribusi

Silakan fork, laporkan issue, atau ajukan pull request untuk perbaikan atau tambahan fitur!

## 📝 Lisensi

MIT License © 2025
