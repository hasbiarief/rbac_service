# Export Swagger to Postman/Insomnia

Quick guide untuk export Swagger documentation ke Postman atau Insomnia.

## 🎯 Quick Start (TL;DR)

**Postman:**
```bash
1. Import file: docs/swagger.json
2. Create environment: base_url = http://localhost:8081
3. Login: POST /api/v1/auth/login
4. Copy token dari response
5. Set environment: token = <copied_token>
6. Test: GET /api/v1/users (token auto-included)
```

**Insomnia:**
```bash
1. Import file: docs/swagger.yaml
2. Create environment: { "base_url": "http://localhost:8081", "token": "" }
3. Login: POST Login with user identity
4. Copy token dari response
5. Update environment: token = <copied_token>
6. Test: GET Get all users (token auto-included)
```

## 📦 Files Available

- `docs/swagger.json` - OpenAPI 3.0 (JSON format)
- `docs/swagger.yaml` - OpenAPI 3.0 (YAML format)
- URL: `http://localhost:8081/swagger/doc.json` (jika server running)

## 🚀 Postman

### Method 1: Import File (Recommended)

1. Buka Postman
2. Click **Import** button (top left)
3. Pilih **File** tab
4. Drag & drop atau browse file `docs/swagger.json`
5. Click **Import**
6. Done! Collection akan muncul di sidebar

### Method 2: Import via URL

1. Pastikan server running: `make dev`
2. Buka Postman → **Import**
3. Pilih **Link** tab
4. Paste URL: `http://localhost:8081/swagger/doc.json`
5. Click **Continue** → **Import**

### Setup Environment

Setelah import, buat environment:

```
base_url: http://localhost:8081
access_token: (akan di-set otomatis setelah login)
```

**Cara set environment:**
1. Click **Environments** (left sidebar)
2. Click **+** untuk create new
3. Nama: `RBAC Development`
4. Add variables:
   - `base_url` = `http://localhost:8081`
   - `access_token` = (kosongkan dulu)
5. Save

### Test Authentication

1. Pilih request **POST /api/v1/auth/login**
2. Body sudah terisi otomatis dari Swagger
3. Click **Send**
4. Copy `access_token` dari response
5. Paste ke environment variable `access_token`
6. Semua request lain akan otomatis gunakan token ini

### Cara Menggunakan Collection (Step-by-Step)

**Setelah import, collection akan terlihat berbeda:**
- Tidak ada folder seperti collection biasa
- Semua endpoints dalam satu list panjang
- Nama request: `POST /api/v1/auth/login` (bukan "Login")

**Langkah-langkah testing:**

1. **Setup Base URL:**
   - Buka **Environments** (left sidebar)
   - Create new environment: `RBAC Dev`
   - Add variable: `baseUrl` = `http://localhost:8081`
   - Save dan select environment ini

2. **Update Request URLs:**
   - Setiap request akan punya URL seperti: `http://localhost:8081/api/v1/auth/login`
   - Ganti dengan: `{{baseUrl}}/api/v1/auth/login`
   - Atau biarkan saja jika tidak mau pakai variable

3. **Test Login:**
   - Cari request: `POST /api/v1/auth/login`
   - Tab **Body** → pilih **raw** → **JSON**
   - Isi body:
   ```json
   {
     "user_identity": "800000001",
     "password": "password123"
   }
   ```
   - Click **Send**
   - Response akan muncul di bawah

4. **Save Token:**
   - Dari response, copy value `data.access_token`
   - Buka **Environments** → edit environment
   - Add variable: `token` = paste token
   - Save

5. **Test Protected Endpoint:**
   - Cari request: `GET /api/v1/users`
   - Tab **Headers** → pastikan ada:
     - Key: `Authorization`
     - Value: `Bearer {{token}}`
   - Jika belum ada, tambahkan manual
   - Click **Send**

**Tips Postman:**
- Gunakan **Search** (Ctrl+K) untuk cari endpoint cepat
- Buat folder manual untuk organize requests
- Duplicate request untuk save different test cases
- Gunakan **Tests** tab untuk auto-save token:
```javascript
// Di request Login, tab Tests:
pm.test("Save token", function () {
    var jsonData = pm.response.json();
    pm.environment.set("token", jsonData.data.access_token);
});
```

## 🔷 Insomnia

### Method 1: Import File (Recommended)

1. Buka Insomnia
2. Click **Create** → **Import From** → **File**
3. Pilih file `docs/swagger.yaml` (Insomnia prefer YAML)
4. Click **Scan** → **Import**
5. Done! Request collection siap digunakan

### Method 2: Import via URL

1. Pastikan server running: `make dev`
2. Buka Insomnia → **Create** → **Import From** → **URL**
3. Paste: `http://localhost:8081/swagger/doc.json`
4. Click **Fetch and Import**

### Setup Environment

1. Click **No Environment** dropdown (top left)
2. Click **Manage Environments**
3. Click **+** untuk create new
4. Nama: `Development`
5. Add JSON:
```json
{
  "base_url": "http://localhost:8081",
  "access_token": ""
}
```
6. Click **Done**

### Test Authentication

1. Pilih request **Login**
2. Click **Send**
3. Copy `access_token` dari response
4. Update environment variable `access_token`
5. Semua request akan otomatis gunakan token

### Cara Menggunakan Collection (Step-by-Step)

**Setelah import, struktur akan berbeda:**
- Requests dikelompokkan berdasarkan Tags (Authentication, Users, dll)
- Nama request lebih deskriptif
- Authorization sudah ter-setup otomatis

**Langkah-langkah testing:**

1. **Setup Environment:**
   - Click dropdown **No Environment** (top left)
   - Select **Manage Environments**
   - Create new: `Development`
   - Add JSON:
   ```json
   {
     "base_url": "http://localhost:8081",
     "token": ""
   }
   ```
   - Click **Done**
   - Select environment `Development`

2. **Update Base URL di Requests:**
   - Setiap request akan punya URL: `http://localhost:8081/api/v1/...`
   - Ganti dengan: `{{ _.base_url }}/api/v1/...`
   - Atau biarkan jika tidak mau pakai variable

3. **Test Login:**
   - Expand folder **Authentication**
   - Click request **POST Login with user identity**
   - Tab **Body** sudah terisi JSON
   - Edit body:
   ```json
   {
     "user_identity": "800000001",
     "password": "password123"
   }
   ```
   - Click **Send**

4. **Save Token Otomatis:**
   - Dari response, copy `data.access_token`
   - Buka **Manage Environments**
   - Edit environment `Development`
   - Update `token` dengan value yang di-copy
   - Save

5. **Test Protected Endpoint:**
   - Expand folder **Users**
   - Click request **GET Get all users**
   - Tab **Auth** → pilih **Bearer Token**
   - Token: `{{ _.token }}`
   - Click **Send**

**Tips Insomnia:**
- Gunakan **Filter** (top right) untuk search requests
- Requests sudah dikelompokkan berdasarkan Tags
- Environment variables: `{{ _.variable_name }}`
- Bisa export/import environment untuk share dengan tim

## 🌐 Online Converters

### Swagger Editor

1. Buka: https://editor.swagger.io/
2. Click **File** → **Import file**
3. Upload `docs/swagger.yaml`
4. Click **Generate Client** → **Postman Collection v2.1**
5. Download dan import ke Postman

### API Transformer

1. Buka: https://www.apimatic.io/transformer
2. Upload `docs/swagger.json`
3. Select output format:
   - Postman Collection v2.1
   - Insomnia v4
   - OpenAPI 3.0
4. Click **Transform**
5. Download hasil konversi

## 💡 Tips

### Postman
- Support JSON dan YAML format
- Auto-generate code snippets (curl, JavaScript, Python, dll)
- Bisa sync collection ke cloud
- Support environment variables dengan `{{variable}}`

### Insomnia
- Prefer YAML format (lebih readable)
- Lightweight dan fast
- Support GraphQL dan gRPC
- Environment variables dengan `{{ _.variable }}`

### Best Practices
1. Gunakan environment variables untuk base_url dan tokens
2. Jangan commit tokens ke git
3. Update collection setiap kali API berubah
4. Gunakan folders untuk organize requests
5. Add tests untuk automated testing

## 🔄 Update Collection

Ketika API berubah:

1. Generate Swagger baru:
```bash
make swagger-gen
```

2. Re-import file ke Postman/Insomnia
3. Pilih **Replace** existing collection
4. Environment variables akan tetap tersimpan

## 🐛 Troubleshooting

**Import gagal:**
- Pastikan file `swagger.json` atau `swagger.yaml` valid
- Cek di Swagger UI dulu: `http://localhost:8081/swagger/index.html`
- Validate di: https://editor.swagger.io/

**Collection terlihat aneh setelah import:**
- Normal! Swagger import berbeda dengan collection biasa
- Postman: semua requests dalam satu list panjang
- Insomnia: requests dikelompokkan berdasarkan Tags
- Solusi: Buat folder manual atau gunakan search

**Request tidak punya body:**
- Setelah import, beberapa request mungkin body-nya kosong
- Buka request → tab **Body** → pilih **raw** → **JSON**
- Copy contoh dari Swagger UI atau dokumentasi
- Atau lihat di `docs/swagger.json` untuk schema

**Authentication tidak work:**
- Pastikan token sudah di-set di environment
- Format harus: `Bearer {{token}}` (ada spasi setelah Bearer)
- Token expire dalam 15 menit, refresh jika perlu
- Cek di tab **Headers**: `Authorization: Bearer {{token}}`

**Base URL salah:**
- Development: `http://localhost:8081`
- Production: sesuaikan dengan domain Anda
- Pastikan tidak ada trailing slash: ❌ `http://localhost:8081/`
- Correct: ✅ `http://localhost:8081`

**Environment variable tidak work:**
- Postman: gunakan `{{variable}}`
- Insomnia: gunakan `{{ _.variable }}`
- Pastikan environment sudah di-select (dropdown top right)
- Variable names case-sensitive

**Response 401 Unauthorized:**
- Token expired (15 menit)
- Token tidak valid
- Token tidak di-set di environment
- Format Authorization header salah
- Solusi: Login ulang dan update token

**Response 404 Not Found:**
- Server tidak running (jalankan: `make dev`)
- Base URL salah
- Endpoint path salah (cek dokumentasi)

**Cara organize collection yang berantakan:**

**Postman:**
1. Buat folder baru: **Authentication**, **Users**, **Companies**, dll
2. Drag & drop requests ke folder yang sesuai
3. Rename requests untuk lebih jelas:
   - `POST /api/v1/auth/login` → `Login`
   - `GET /api/v1/users` → `Get All Users`

**Insomnia:**
- Requests sudah dikelompokkan otomatis berdasarkan Tags
- Tidak perlu organize manual

## 📚 Resources

- Postman Docs: https://learning.postman.com/docs/getting-started/importing-and-exporting-data/
- Insomnia Docs: https://docs.insomnia.rest/insomnia/import-export-data
- OpenAPI Spec: https://swagger.io/specification/
- Swagger Editor: https://editor.swagger.io/
