# Auth Service

Microservice untuk autentikasi menggunakan Golang, Gin, GORM, PostgreSQL, dan JWT (RSA256 dengan standar JWKS). Layanan ini juga menyediakan dokumentasi API interaktif menggunakan antarmuka modern dari Scalar.

## Fitur Utama

- **Registrasi & Login**: Menggunakan hashing password dengan `bcrypt`.
- **JWT (RSA256)**: Token ditandatangani menggunakan algoritma asimetris RSA untuk keamanan maksimal.
- **JWKS (JSON Web Key Set)**: Mengekspos endpoint `/.well-known/jwks.json` standar agar layanan eksternal dapat memvalidasi token JWT secara mandiri (production-grade).
- **Refresh Token**: Mekanisme aman untuk memperbarui access token yang kedaluwarsa.
- **API Documentation**: Dokumentasi otomatis yang di-generate menggunakan `swaggo` dan ditampilkan menggunakan `gin-openapi` (Scalar UI).

## Persiapan Awal (Setup)

1. Pastikan Anda memiliki PostgreSQL yang sedang berjalan.
2. Buat database baru di PostgreSQL dengan nama `ms-auth-service`.
3. Pastikan Anda memiliki file konfigurasi `.env` di **dalam folder `auth`** (atau sesuaikan pengaturannya jika dijalankan dari root):
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=f0bima
   DB_PASSWORD=f0bima
   DB_NAME=ms-auth-service

   JWT_SECRET=supersecret
   JWT_EXPIRE=15m
   REFRESH_EXPIRE=168h
   ```

## Menjalankan Migrasi Database

Layanan ini menggunakan [golang-migrate](https://github.com/golang-migrate/migrate) untuk mengatur skema database.
Skrip `migrate.bat` disediakan agar Anda bisa menjalankan migrasi dengan mudah di Windows. Skrip ini akan secara otomatis membaca variabel dari file `.env` di folder saat ini.

- **Menerapkan semua migrasi (Up)**:
  ```bash
  .\migrate.bat up
  ```
- **Membatalkan migrasi (Down/Rollback)**:
  ```bash
  .\migrate.bat down
  ```
- **Membuat file migrasi baru**:
  ```bash
  .\migrate.bat create nama_fitur_baru
  ```

_(Catatan: Anda juga bisa menggunakan Makefile jika memiliki `make` terinstal, misal: `make migrate-up`)_

## Dokumentasi API (Swagger & Scalar)

Aplikasi ini menggunakan anotasi _Swagger_ pada kode Go (seperti di `main.go` dan `handler.go`) untuk mendeskripsikan API. Scalar UI (di endpoint `/docs/`) akan merender dokumentasi ini secara visual.

**PENTING**: Scalar UI tidak otomatis membaca perubahan pada baris komentar Anda. Scalar hanya membaca file statis `docs/swagger.json`.

Oleh karena itu, **setiap kali Anda menambahkan endpoint baru atau mengubah komentar (anotasi) pada kode**, Anda **wajib** menjalankan perintah berikut untuk memperbarui file dokumentasi:

```bash
# Instal swag CLI (jika belum pernah)
go install github.com/swaggo/swag/cmd/swag@latest

# Jalankan perintah ini di dalam folder auth/
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -o ./docs
```

_(Setelah perintah di atas berhasil dijalankan, Anda bisa me-refresh Scalar UI di peramban Anda untuk melihat perubahannya.)_

## Live Reloading dengan Air

Project ini mendukung **Air** untuk _live-reload_ saat masa pengembangan. Jika Anda belum menginstal Air, ikuti langkah penting berikut:

1. **Install Air**:
   Versi terbaru menggunakan repository `air-verse`:

   ```bash
   go install github.com/air-verse/air@latest
   ```

   Pastikan folder `$(go env GOPATH)/bin` atau `$(go env GOBIN)` sudah masuk ke `PATH` sistem Anda agar perintah `air` bisa dikenali.

2. **Konfigurasi**:
   Konfigurasi sudah diatur di dalam `.air.toml` (melakukan build dari `cmd/api` ke folder `tmp/`). Anda tidak perlu melakukan `air init` lagi.

## Menjalankan Server

1. Pastikan folder `keys/` memiliki file `private.pem` dan `public.pem`. (Kunci RSA 2048-bit).
2. Jalankan aplikasi secara manual:
   ```bash
   go run cmd/api/main.go
   ```
   Atau jalankan dengan **Air** agar server otomatis di-_restart_ setiap kali Anda mengubah file `.go`:
   ```bash
   air
   ```

Setelah server berjalan, Anda dapat mengakses:

- **API Endpoint Utama**: `http://localhost:8080/auth/...`
- **JWKS Endpoint**: `http://localhost:8080/.well-known/jwks.json`
- **Dokumentasi API (Scalar UI)**: `http://localhost:8080/docs/`
