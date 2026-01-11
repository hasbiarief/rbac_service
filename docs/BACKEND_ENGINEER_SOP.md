# Backend Engineer SOP — Ringkas
Tujuan: Panduan singkat untuk setup, pengembangan fitur, dan deploy service RBAC.

Quick start:
1. git clone <repo> && cd rbac-service
2. go mod download
3. cp .env.example .env (isi DB_*, REDIS_ADDR, JWT_SECRET)
4. createdb huminor_rbac && make migrate-up
5. air / go run cmd/api/main.go

Workflow singkat:
- Desain API & DTO → Tambah migration (jika perlu) → Repo → Service → Handler → Routes → Test
- Gunakan DTO, interfaces, mapper, validation middleware, dan constants.

