Viewed README.md:1-14
Listed directory rules
Viewed devops-ci.md:1-366
Viewed microservices-architecture.md:1-480

Terkait **Observability** (Keteramatan Sistem), _codebase_ `spd-web-v2` memiliki standar _Production-Grade_ yang sangat komprehensif menggunakan pilar utama observability: **Traces (Tempo), Logs (Loki), dan Metrics (Prometheus)** yang semuanya diikat oleh standar **OpenTelemetry (OTel)**.

Berikut adalah ekstraksi aturan _Production-Grade_ terkait Observability yang bersifat agnostik dan wajib diterapkan di AOS:

### 1. Distributed Tracing (OpenTelemetry & Tempo)

TraceID adalah "benang merah" yang mengikat semua log dan _request_ dari ujung depan (Frontend) sampai ujung belakang (Database/Microservices).

- **Aturan Implementasi:**
  - **Wajib Span Wrapping:** Semua fungsi _use-case_ atau pemanggilan API eksternal WAJIB dibungkus dengan _Span_ (`tracer.startActiveSpan`).
  - **Status & Exception Recording:** Jika terjadi _error_ di dalam blok kode, sistem tidak boleh hanya melempar _throw error_, tapi wajib mencatatnya di _span_ (`span.setStatus(ERROR)` dan `span.recordException(error)`).
  - **Context Propagation (Sangat Krusial):** Jika _Service A_ memanggil _Service B_, _Service A_ **WAJIB** meneruskan _TraceID_ melalui _HTTP Header_ (biasanya header `traceparent` atau `x-trace-id`). Tanpa ini, _trace_ akan terputus.

### 2. Structured Logging (Loki)

Di level produksi, pesan log berupa teks biasa (seperti `console.log("data berhasil disimpan")`) **DILARANG KERAS**.

- **Aturan Implementasi:**
  - **Wajib berformat JSON (Structured):** Logger harus mencetak log dalam format objek (JSON).
  - **Wajib menyertakan Konteks:** Setiap _log_ WAJIB memasukkan `traceId` (agar bisa dicocokkan dengan Tempo di Grafana), `userId`, `event_name`, dan `timestamp`.
  - **Contoh Penerapan Universal:**
    ```json
    {
      "level": "info",
      "event": "trip_plan_created",
      "userId": "usr-123",
      "tripPlanId": "tp-456",
      "traceId": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
      "message": "Trip plan successfully created"
    }
    ```
  - **Leveling yang Ketat:** Gunakan `error` HANYA untuk kegagalan sistem, `warn` untuk potensi bahaya/depresiasi, `info` untuk _business action_ yang sukses (seperti persetujuan, pembayaran), dan `debug` untuk _tracing_ data internal.

### 3. Monitoring, Metrics, & Alerting (Prometheus & Grafana)

_Metrics_ digunakan untuk melihat "kesehatan" aplikasi secara makro.

- **Aturan Implementasi:**
  - **Golden Signals Monitoring:** Sistem wajib mengukur _Response Time/Latency_ (p50, p95, p99), _Error Rates_ (rasio HTTP 5xx & 4xx), dan Saturasi (_CPU/Memory Usage_, _Database Connection Pool_).
  - **Alerting Thresholds:** Sistem harus memicu peringatan (ke Slack/Email) jika:
    - _Error rate_ > 1%
    - _Response time_ persentil-99 (p99) > 2 detik
    - Pemakaian memori > 90%
    - Koneksi database habis (_pool exhausted_).

### 4. Health & Readiness Probes (Kubernetes Ready)

Sistem produksi yang agnostik harus memiliki gerbang pemeriksaan kesehatan agar _Load Balancer_ atau _Kubernetes_ tahu aplikasi sedang hidup.

- **Aturan Implementasi:**
  - Wajib mengekspos rute HTTP `/api/health`.
  - **Liveness vs Readiness:** Pemeriksaan tidak boleh hanya mereturn HTTP 200 "OK". Sistem wajib melakukan tes _ping_ aktif ke _Database_, _Redis_, atau layanan pihak ketiga yang menjadi dependensi utamanya. Jika DB mati, _healthcheck_ wajib mereturn `HTTP 503 Service Unavailable`.

### Keterikatan di Grafana (The Magic of OTel)

Mengapa standar ini sangat krusial di Production?
Karena dengan standar di atas, ketika Anda melihat lonjakan _Error 500_ di dasbor metrik (**Prometheus**), Anda bisa mengklik titik lonjakan tersebut. Grafana akan langsung membuka **Tempo** (Tracing) untuk menunjukkan fungsi mana yang lambat atau _crash_. Dan dari _span_ Tempo tersebut, Anda akan menemukan satu **TraceID**. Dengan mengklik TraceID tersebut, **Loki** akan memunculkan semua baris **Log** JSON yang dieksekusi selama _request_ tersebut secara spesifik!

Inilah yang dinamakan standar Observability kelas produksi—bukan sekadar mencetak _error_ ke terminal, tetapi membangun ekosistem di mana metrik, _trace_, dan log saling berbicara. Aturan ini sangat mudah diwujudkan di AOS sebagai **Global Governance Rules** (`observability.yaml`).
