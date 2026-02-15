# Postman to Swagger Migration Script

Script untuk memigrasikan koleksi Postman ke anotasi Swagger secara otomatis.

## Fitur

- ✅ Konversi koleksi Postman ke anotasi Swagger
- ✅ Generasi file anotasi per modul
- ✅ Validasi hasil konversi
- ✅ Generasi dokumentasi Swagger otomatis
- ✅ Laporan migrasi lengkap

## Prasyarat

- Go 1.21 atau lebih tinggi
- swag CLI tool (`go install github.com/swaggo/swag/cmd/swag@latest`)
- jq (untuk parsing JSON)

## Penggunaan

### Menggunakan Make

```bash
# Dengan file default (postman_collection.json)
make migrate-postman

# Dengan file custom
make migrate-postman ARGS="-f path/to/collection.json"

# Dengan output directory custom
make migrate-postman ARGS="-f collection.json -o /path/to/output"
```

### Menggunakan Script Langsung

```bash
# Dengan file default
./scripts/migrate-postman.sh

# Dengan file custom
./scripts/migrate-postman.sh -f path/to/collection.json

# Dengan opsi lengkap
./scripts/migrate-postman.sh -f collection.json -o . -d docs
```

### Menggunakan Environment Variables

```bash
POSTMAN_FILE=my-collection.json ./scripts/migrate-postman.sh
```

## Opsi

| Opsi | Deskripsi | Default |
|------|-----------|---------|
| `-f, --file` | File koleksi Postman | `postman_collection.json` |
| `-o, --output` | Direktori output | `.` (current directory) |
| `-d, --docs` | Direktori dokumentasi | `docs` |
| `-h, --help` | Tampilkan bantuan | - |

## Environment Variables

| Variable | Deskripsi | Default |
|----------|-----------|---------|
| `POSTMAN_FILE` | Path ke file koleksi Postman | `postman_collection.json` |
| `OUTPUT_DIR` | Direktori output | `.` |
| `DOCS_DIR` | Direktori dokumentasi | `docs` |
| `SWAGGER_CMD` | Path ke swagger CLI tool | `./bin/swagger` |

## Proses Migrasi

Script akan melakukan langkah-langkah berikut:

1. **Validasi Prasyarat**
   - Memeriksa instalasi Go
   - Memeriksa instalasi swag
   - Membangun swagger CLI tool jika belum ada

2. **Validasi File Postman**
   - Memeriksa keberadaan file
   - Memvalidasi format JSON
   - Memverifikasi struktur koleksi Postman

3. **Konversi**
   - Mengelompokkan endpoint berdasarkan modul
   - Mengonversi setiap endpoint ke anotasi Swagger
   - Menghasilkan file `swagger.go` untuk setiap modul

4. **Generasi Dokumentasi**
   - Menjalankan `swag init`
   - Menghasilkan `swagger.json` dan `swagger.yaml`
   - Menghasilkan `docs.go` untuk embedding

5. **Validasi Hasil**
   - Memeriksa file anotasi yang dihasilkan
   - Memvalidasi dokumentasi yang dihasilkan
   - Menjalankan validator anotasi

6. **Laporan**
   - Menghasilkan laporan migrasi
   - Menampilkan statistik konversi
   - Memberikan langkah selanjutnya

## Output

Script akan menghasilkan:

- **File Anotasi**: `internal/modules/{module}/docs/swagger.go`
- **Dokumentasi**: `docs/swagger.json`, `docs/swagger.yaml`, `docs/docs.go`
- **Laporan**: `migration-report.txt`

## Contoh Output

```
========================================
Postman to Swagger Migration Tool
========================================

========================================
Checking Requirements
========================================
✓ Go is installed
✓ swag is installed
✓ Swagger CLI tool found

========================================
Validating Postman Collection
========================================
✓ Valid Postman collection: My API Collection
ℹ Found 25 endpoints

========================================
Converting Postman Collection to Swagger Annotations
========================================
ℹ Running converter...
✓ Conversion completed
ℹ Total endpoints: 25
ℹ Converted: 23
⚠ Failed: 2

⚠ Failed endpoints:
  - Complex Endpoint: unsupported request format
  - Legacy Endpoint: missing required fields

✓ Generated 5 annotation files:
  - internal/modules/auth/docs/swagger.go
  - internal/modules/user/docs/swagger.go
  - internal/modules/company/docs/swagger.go
  - internal/modules/branch/docs/swagger.go
  - internal/modules/role/docs/swagger.go

========================================
Generating Swagger Documentation
========================================
ℹ Running swag init...
✓ Swagger documentation generated
✓ Generated: docs/swagger.json
ℹ Documented endpoints: 23
✓ Generated: docs/swagger.yaml

========================================
Validating Conversion Results
========================================
ℹ Checking annotation files...
✓ Found 5 annotation files
ℹ Validating generated documentation...
✓ swagger.json is valid
✓ Found 23 paths in documentation
ℹ Running annotation validator...
✓ All annotations are valid

========================================
Migration Report
========================================
Postman to Swagger Migration Report
====================================

Date: 2024-01-15 10:30:00
Postman Collection: my-collection.json
Collection Name: My API Collection

Results:
--------
Annotation files generated: 5
Endpoints documented: 23

Generated Files:
----------------
internal/modules/auth/docs/swagger.go
internal/modules/user/docs/swagger.go
internal/modules/company/docs/swagger.go
internal/modules/branch/docs/swagger.go
internal/modules/role/docs/swagger.go

docs/swagger.json
docs/swagger.yaml
docs/docs.go

Next Steps:
-----------
1. Review generated annotation files in internal/modules/*/docs/swagger.go
2. Update annotations with accurate type information and descriptions
3. Run 'make swagger-gen' to regenerate documentation after changes
4. Test Swagger UI at http://localhost:8081/swagger/index.html
5. Validate all endpoints work correctly

✓ Report saved to: migration-report.txt

========================================
Migration Completed Successfully!
========================================
✓ Your Postman collection has been migrated to Swagger annotations
ℹ Review the migration report for next steps
```

## Troubleshooting

### Error: "swag is not installed"

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Error: "Failed to build swagger CLI tool"

Pastikan file `cmd/swagger/main.go` ada dan dapat dikompilasi:

```bash
go build -o bin/swagger cmd/swagger/main.go
```

### Error: "Invalid JSON in Postman collection file"

Validasi file Postman Anda:

```bash
jq empty your-collection.json
```

### Konversi Gagal untuk Beberapa Endpoint

Periksa laporan migrasi untuk detail endpoint yang gagal. Anda mungkin perlu:
- Memperbaiki format request di Postman
- Menambahkan informasi yang hilang
- Mengonversi endpoint tersebut secara manual

## Langkah Setelah Migrasi

1. **Review Anotasi**
   - Periksa file `internal/modules/*/docs/swagger.go`
   - Perbarui tipe data yang tidak akurat
   - Tambahkan deskripsi yang lebih detail

2. **Update Tipe Data**
   - Ganti `object` dengan tipe struct yang sesuai
   - Tambahkan referensi ke DTO yang benar
   - Contoh: `object` → `auth.LoginRequest`

3. **Regenerasi Dokumentasi**
   ```bash
   make swagger-gen
   ```

4. **Test Swagger UI**
   - Jalankan aplikasi: `make run`
   - Buka: http://localhost:8081/swagger/index.html
   - Test setiap endpoint

5. **Validasi**
   ```bash
   make swagger-validate
   ```

## Tips

- Gunakan nama yang konsisten untuk modul di Postman (folder structure)
- Tambahkan deskripsi yang jelas di Postman sebelum migrasi
- Kelompokkan endpoint berdasarkan modul untuk hasil yang lebih baik
- Review dan update anotasi setelah migrasi untuk akurasi maksimal

## Lihat Juga

- [Swagger Documentation](../docs/swagger.json)
- [swaggo/swag Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
