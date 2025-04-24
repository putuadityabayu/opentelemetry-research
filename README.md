# 📊 opentelemetry-research

Eksperimen observabilitas menggunakan OpenTelemetry di aplikasi Go, lengkap dengan integrasi Prometheus dan Grafana. Proyek ini mensimulasikan perekaman metrik seperti nilai acak dan penggunaan memori heap secara real-time.

## ✨ Fitur

- Integrasi OpenTelemetry SDK di aplikasi Go
- Ekspor metrik ke Prometheus melalui OpenTelemetry Collector
- Visualisasi metrik menggunakan Grafana
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
- OpenTelemetry Collector (terhubung internal container)

### 3. Jalankan Aplikasi Go

```bash
go run main.go
```

## 📈 Akses dan Visualisasi

- Grafana: http://localhost:3000
  - Login: admin / admin
  - Datasource Prometheus akan otomatis terprovisi.
- Prometheus UI: http://localhost:9090
  - Gunakan untuk query metrik seperti `random_value` atau `memory_heap`.

## 📊 Metrik yang Disediakan

- `random_value` (Gauge - async): Nilai acak antara 0 dan 100 yang diperbarui setiap polling.
- `memory_heap` (Gauge - async): Heap memory dari aplikasi Go (dalam byte), diambil dari `runtime.ReadMemStats()`.

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

exporters:
  prometheus:
    endpoint: "0.0.0.0:8889"

service:
  pipelines:
    metrics:
      receivers: [otlp]
      exporters: [prometheus]

```

### `prometheus/prometheus.yml`

Konfigurasi untuk prometheus:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'otel-collector'
    static_configs:
      - targets: ['otel-collector:8889']  # scrape metrik dari Collector
```

### `grafana/provisioning/datasources/datasource.yml`

Konfigurasi otomatis datasource Grafana:

```yaml
apiVersion: 1
deleteDatasources:
  - name: 'Prometheus'
    orgId: 1
prune: true

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    editable: false
```

## 🛠️ Dependencies

Pastikan Anda telah menginstall:
- [Docker](https://www.docker.com/)
- [Go](https://golang.org/) (versi terbaru)

## 📚 Referensi

- [OpenTelemetry Go SDK](https://opentelemetry.io/docs/instrumentation/go/)
- [Prometheus](https://prometheus.io/)
- [Grafana](https://grafana.com/)

## 🤝 Kontribusi

Silakan fork, laporkan issue, atau ajukan pull request untuk perbaikan atau tambahan fitur!

## 📝 Lisensi

MIT License © 2025
