# Air Live Reload - Panduan Penggunaan

## Apa itu Air?

Air adalah tool untuk live reload development server Go. Ketika Anda mengubah kode, Air akan otomatis rebuild dan restart server tanpa perlu manual restart.

## Instalasi Air

```bash
# Install Air secara global
go install github.com/cosmtrek/air@latest

# Pastikan $GOPATH/bin ada di PATH Anda
export PATH=$PATH:$(go env GOPATH)/bin
```

## Cara Menjalankan dengan Air

### Opsi 1: Jalankan dengan konfigurasi default
```bash
# Di root directory project
air
```

### Opsi 2: Jalankan dengan verbose mode (untuk debugging)
```bash
air -v
```

### Opsi 3: Jalankan dengan konfigurasi custom
```bash
air -c .air.toml
```

## Konfigurasi Air (.air.toml)

Project ini sudah memiliki konfigurasi Air di file `.air.toml`:

```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main ./cmd/api"
  bin = "./tmp/main"
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "node_modules", "bin"]
  exclude_regex = ["_test.go"]
  delay = 1000
```

### Penjelasan Konfigurasi:
- `cmd`: Command untuk build aplikasi
- `bin`: Path ke binary yang dihasilkan
- `include_ext`: File extension yang akan di-watch
- `exclude_dir`: Directory yang diabaikan
- `exclude_regex`: Pattern file yang diabaikan
- `delay`: Delay sebelum restart (ms)

## Workflow Development dengan Air

1. **Start development server:**
   ```bash
   air
   ```

2. **Edit kode** - Air akan otomatis detect perubahan

3. **Air akan:**
   - Compile ulang kode
   - Restart server
   - Menampilkan log build

4. **Test API** - Server langsung siap digunakan

## Output Air

Ketika menjalankan `air`, Anda akan melihat output seperti:
```
  __    _   ___  
 / /\  | | | |_) 
/_/--\ |_| |_| \_ v1.63.4, built with Go go1.25.2

watching .
watching cmd
watching cmd/api
watching config
watching internal
...
building...
running...
```

## Tips Penggunaan

### 1. Monitoring Changes
Air akan menampilkan file mana yang berubah:
```
main.go has changed
building...
running...
```

### 2. Build Errors
Jika ada error kompilasi, Air akan menampilkan error dan tidak restart:
```
building...
# gin-scalable-api/internal/handlers
internal/handlers/user_handler.go:25:2: syntax error
failed to build, error: exit status 2
```

### 3. Exclude Files
Untuk mengabaikan file tertentu, edit `.air.toml`:
```toml
exclude_regex = ["_test.go", ".*_temp.go"]
```

### 4. Custom Build Command
Untuk build command khusus:
```toml
[build]
  cmd = "go build -tags dev -o ./tmp/main ./cmd/api"
```

## Troubleshooting

### Air tidak terinstall
```bash
# Check instalasi
air --version

# Jika tidak ada, install ulang
go install github.com/cosmtrek/air@latest
```

### Port sudah digunakan
```bash
# Kill process yang menggunakan port 8081
lsof -ti:8081 | xargs kill -9

# Atau ubah port di .env
PORT=8082
```

### Air tidak detect perubahan
```bash
# Restart Air
Ctrl+C
air
```

### Build terlalu lambat
Edit `.air.toml` untuk exclude lebih banyak directory:
```toml
exclude_dir = ["assets", "tmp", "vendor", "testdata", "node_modules", "bin", "docs", "migrations"]
```

## Perbandingan dengan Manual Run

| Aspek | Manual Run | Air Live Reload |
|-------|------------|-----------------|
| **Restart** | Manual `Ctrl+C` + `./server` | Otomatis |
| **Build Time** | Sama | Sama |
| **Development Speed** | Lambat | Cepat |
| **Memory Usage** | Rendah | Sedikit lebih tinggi |
| **Debugging** | Mudah | Mudah |

## Best Practices

1. **Gunakan Air untuk development** - Jangan untuk production
2. **Exclude test files** - Untuk performa lebih baik
3. **Monitor build errors** - Air akan menampilkan error dengan jelas
4. **Restart Air** jika ada masalah konfigurasi
5. **Gunakan verbose mode** (`air -v`) untuk debugging

## Alternatif Air

Jika Air bermasalah, Anda bisa menggunakan:

### 1. Manual Watch Script
```bash
# watch.sh
#!/bin/bash
while true; do
  go build -o server cmd/api/main.go && ./server &
  PID=$!
  inotifywait -e modify -r . --exclude '(\.git|tmp|vendor)'
  kill $PID
done
```

### 2. Makefile dengan Watch
```makefile
watch:
	@while true; do \
		make build && ./server & \
		PID=$$!; \
		inotifywait -e modify -r . --exclude '(\.git|tmp|vendor)'; \
		kill $$PID; \
	done
```

### 3. VS Code Tasks
Buat `.vscode/tasks.json`:
```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Go: Build and Run",
      "type": "shell",
      "command": "go build -o server cmd/api/main.go && ./server",
      "group": "build",
      "presentation": {
        "echo": true,
        "reveal": "always",
        "focus": false,
        "panel": "new"
      }
    }
  ]
}
```

---

**Catatan:** Air sangat membantu untuk development karena menghemat waktu restart manual. Gunakan selalu untuk development, tapi jangan untuk production deployment.