# 📦 Rekomendasi Penyederhanaan Struktur Project Gin (Express-style)

> Konteks:
> - Framework: **Gin**
> - Bahasa: **Go**
> - ORM: **Tidak pakai GORM (raw SQL + file migration SQL)**
> - Background developer: **TypeScript / Express**
> - Domain: **RBAC Service**

---

## 🎯 Masalah Utama Struktur Saat Ini

1. Struktur terlalu **horizontal (per layer)**  
   - handlers / service / repository / dto / validation terpisah
   - 1 endpoint → buka banyak folder

2. Banyak folder **over-engineering**
   - `interfaces`
   - `mapper`
   - `dto` global

3. Folder `pkg` berisi **domain logic**
   - `pkg/rbac`
   - `pkg/model`
   ➜ seharusnya jadi module

4. Mental model tidak selaras dengan **Express**
   - Express → route → validation → controller → model
   - Gin bisa (dan sebaiknya) meniru alur ini

---

## ✅ Prinsip Desain yang Direkomendasikan

- **1 fitur = 1 folder**
- Struktur **vertikal (per module)**, bukan horizontal
- Interface **dekat dengan pemakai**, bukan di folder khusus
- `pkg` hanya untuk **kode generik & reusable**
- Raw SQL → repository **simple & eksplisit**
- Validation **dekat dengan route / handler**

---

## 🗂️ Struktur Folder yang Direkomendasikan

rbac-service/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── migrate/
│       └── main.go

├── config/
│   └── config.go

├── internal/
│   ├── app/
│   │   └── server.go
│
│   ├── middleware/
│   │   ├── auth.go
│   │   └── rbac.go
│
│   ├── modules/
│   │   ├── auth/
│   │   ├── user/
│   │   │   ├── route.go        # route
│   │   │   ├── handler.go      # controller
│   │   │   ├── service.go      # business logic
│   │   │   ├── repository.go   # raw SQL
│   │   │   ├── model.go        # struct db
│   │   │   └── validator.go    # request validation
│   │   │
│   │   ├── role/
│   │   ├── permission/
│   │   └── policy/
│
│   └── shared/
│       ├── db/
│       ├── response/
│       ├── errors/
│       ├── pagination/
│       └── query/
│
├── migrations/        # pure SQL
│   ├── 001_init.sql
│   └── 002_rbac.sql
│
├── scripts/
├── docs/
└── tmp/


---

## 🔁 Mapping Struktur Lama → Baru

### 🔴 Folder yang Sebaiknya Dihapus

| Folder Lama | Alasan |
|------------|-------|
| `internal/interfaces` | Go tidak butuh folder interface |
| `internal/mapper` | Mapping bisa inline / function kecil |
| `internal/dto` (global) | DTO sebaiknya per module |

---

### 🟡 Folder yang Dipindahkan ke Module

| Folder Lama | Tujuan Baru |
|------------|------------|
| `internal/handlers` | `modules/*/handler.go` |
| `internal/service` | `modules/*/service.go` |
| `internal/repository` | `modules/*/repository.go` |
| `internal/validation` | `modules/*/validator.go` |
| `internal/models` | `modules/*/model.go` |

---

### 🟢 `pkg` yang Tetap Dipertahankan

Gunakan `pkg` **hanya untuk reusable & generic**

pkg/
├── logger
├── password
├── token
├── ratelimiter
├── utils


❌ Jangan taruh domain seperti `rbac`, `user`, `role` di `pkg`

---

## 🔧 Contoh Repository (Raw SQL, Tanpa GORM)

```go
func (r *Repository) FindByID(ctx context.Context, id int64) (*User, error) {
    row := r.db.QueryRowContext(
        ctx,
        `SELECT id, email FROM users WHERE id = $1`,
        id,
    )

    var u User
    if err := row.Scan(&u.ID, &u.Email); err != nil {
        return nil, err
    }

    return &u, nil
}

🧭 Rule of Thumb (Wajib Diingat)
Kalau nambah 1 fitur tapi harus buka > 2 folder → struktur terlalu ribet
Interface hanya dibuat kalau memang dibutuhkan
Module RBAC bukan middleware, middleware hanya enforcement
Raw SQL + Go struct sudah cukup untuk 90% use case

🏁 Kesimpulan
Struktur awal cocok untuk enterprise Java-style
Untuk Gin + Express mindset → module-based structure lebih optimal
Struktur baru:
Lebih cepat dikembangkan
Lebih mudah dipahami
Siap diskalakan ke microservice