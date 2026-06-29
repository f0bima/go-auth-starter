module.exports = {
  '*.go': [
    'gofmt -w',
    // Gunakan fungsi untuk mencegah lint-staged menyisipkan nama file di akhir perintah.
    // golangci-lint bekerja lebih baik jika dijalankan pada level direktori/proyek,
    // bukan pada daftar file spesifik dari berbagai package yang berbeda.
    () => 'golangci-lint run --timeout=5m',
  ],
  '*.{js,ts,json,md,yml,yaml}': ['prettier --write'],
}
