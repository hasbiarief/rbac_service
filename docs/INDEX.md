# Documentation Index - RBAC Service

## 📚 Dokumentasi Utama

### 🚀 [Quick Start Guide](QUICK_START.md)
Panduan setup cepat untuk development dan production. Termasuk Makefile commands.

### 🏗️ [Project Structure](PROJECT_STRUCTURE.md)
Penjelasan lengkap tentang module-based architecture, struktur folder, dan prinsip desain.

### 👨‍💻 [Backend Engineer Rules](ENGINEER_RULES.md)
Panduan lengkap untuk backend engineer: development workflow, module structure, best practices.

### 📡 [API Documentation Module](apidoc/)
Dokumentasi lengkap API Documentation System dan semua API endpoints.

### 🔗 [Integration Guides](integration/)
Panduan integrasi authentication service untuk external applications.

## 🧪 Testing

### Postman Collection
- **Collection**: `HUMINOR_RBAC_API_MODULE_BASED.postman_collection.json`
- **Environment**: `HUMINOR_RBAC_Environment_Module_Based.postman_environment.json`

Import ke Postman untuk testing lengkap semua endpoints.

## 🎯 Quick Links

**Untuk memulai development:**
1. Baca [Quick Start Guide](QUICK_START.md)
2. Pahami [Project Structure](PROJECT_STRUCTURE.md)
3. Ikuti [Backend Engineer Rules](ENGINEER_RULES.md)

**Untuk testing API:**
1. Import Postman collection
2. Baca [API Documentation](apidoc/)
3. Test dengan Postman

## 📖 Urutan Belajar

1. **Setup** → [Quick Start Guide](QUICK_START.md)
2. **Arsitektur** → [Project Structure](PROJECT_STRUCTURE.md)
3. **Development** → [Backend Engineer Rules](ENGINEER_RULES.md)
4. **API** → [API Documentation Module](apidoc/)
5. **Testing** → Postman Collection

## 🔑 Key Concepts

- **Module-Based Architecture**: 1 fitur = 1 folder (Express.js style)
- **7 Files per Module**: route, handler, service, repository, model, dto, validator
- **No Cross-Module Imports**: Setiap module independent
- **Local Models**: Tidak ada shared models
- **Raw SQL**: Tanpa ORM untuk performa optimal

## 📝 Project Info

- **Framework**: Gin (Go)
- **Database**: PostgreSQL (raw SQL)
- **Cache**: Redis
- **Auth**: JWT with refresh token
- **Architecture**: Module-Based (vertical structure)
