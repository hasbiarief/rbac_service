# Module Structure Refactoring Plan

## 🎯 Tujuan
Menyederhanakan struktur modules dari 7 files menjadi 5 files per module untuk meningkatkan developer experience dan mengurangi cognitive load.

## 📋 Perubahan Structure

### Sebelum (7 files):
```
modules/{module}/
├── dto.go
├── handler.go      ← akan digabung ke route.go
├── model.go
├── repository.go
├── route.go        ← akan menerima handler logic
├── service.go
└── validator.go    ← akan digabung ke dto.go
```

### Sesudah (5 files):
```
modules/{module}/
├── dto.go          ← + validation logic dari validator.go
├── model.go
├── repository.go
├── route.go        ← + handler logic dari handler.go
└── service.go
```

## 🔄 Refactoring Steps

### ✅ Module 1: /audit - COMPLETED ✅
- [x] **Step 1.1**: Merge `validator.go` content ke `dto.go`
- [x] **Step 1.2**: Merge `handler.go` content ke `route.go`
- [x] **Step 1.3**: Update imports di `route.go` jika diperlukan
- [x] **Step 1.4**: Delete `validator.go` dan `handler.go`
- [x] **Step 1.5**: Test module audit masih berfungsi normal

### ✅ Module 2: /auth
- [ ] **Step 2.1**: Merge `validator.go` content ke `dto.go`
- [ ] **Step 2.2**: Merge `handler.go` content ke `route.go`
- [ ] **Step 2.3**: Update imports di `route.go` jika diperlukan
- [ ] **Step 2.4**: Delete `validator.go` dan `handler.go`
- [ ] **Step 2.5**: Test module auth masih berfungsi normal

### ✅ Module 3: /branch
- [ ] **Step 3.1**: Merge `validator.go` content ke `dto.go`
- [ ] **Step 3.2**: Merge `handler.go` content ke `route.go`
- [ ] **Step 3.3**: Update imports di `route.go` jika diperlukan
- [ ] **Step 3.4**: Delete `validator.go` dan `handler.go`
- [ ] **Step 3.5**: Test module branch masih berfungsi normal

### ✅ Module 4: /company
- [ ] **Step 4.1**: Merge `validator.go` content ke `dto.go`
- [ ] **Step 4.2**: Merge `handler.go` content ke `route.go`
- [ ] **Step 4.3**: Update imports di `route.go` jika diperlukan
- [ ] **Step 4.4**: Delete `validator.go` dan `handler.go`
- [ ] **Step 4.5**: Test module company masih berfungsi normal

### ✅ Module 5: /module
- [ ] **Step 5.1**: Merge `validator.go` content ke `dto.go`
- [ ] **Step 5.2**: Merge `handler.go` content ke `route.go`
- [ ] **Step 5.3**: Update imports di `route.go` jika diperlukan
- [ ] **Step 5.4**: Delete `validator.go` dan `handler.go`
- [ ] **Step 5.5**: Test module module masih berfungsi normal

### ✅ Module 6: /role
- [ ] **Step 6.1**: Merge `validator.go` content ke `dto.go`
- [ ] **Step 6.2**: Merge `handler.go` content ke `route.go`
- [ ] **Step 6.3**: Update imports di `route.go` jika diperlukan
- [ ] **Step 6.4**: Delete `validator.go` dan `handler.go`
- [ ] **Step 6.5**: Test module role masih berfungsi normal

### ✅ Module 7: /subscription
- [ ] **Step 7.1**: Merge `validator.go` content ke `dto.go`
- [ ] **Step 7.2**: Merge `handler.go` content ke `route.go`
- [ ] **Step 7.3**: Update imports di `route.go` jika diperlukan
- [ ] **Step 7.4**: Delete `validator.go` dan `handler.go`
- [ ] **Step 7.5**: Test module subscription masih berfungsi normal

### ✅ Module 8: /unit
- [ ] **Step 8.1**: Merge `validator.go` content ke `dto.go`
- [ ] **Step 8.2**: Merge `handler.go` content ke `route.go`
- [ ] **Step 8.3**: Update imports di `route.go` jika diperlukan
- [ ] **Step 8.4**: Delete `validator.go` dan `handler.go`
- [ ] **Step 8.5**: Test module unit masih berfungsi normal

### ✅ Module 9: /user
- [ ] **Step 9.1**: Merge `validator.go` content ke `dto.go`
- [ ] **Step 9.2**: Merge `handler.go` content ke `route.go`
- [ ] **Step 9.3**: Update imports di `route.go` jika diperlukan
- [ ] **Step 9.4**: Delete `validator.go` dan `handler.go`
- [ ] **Step 9.5**: Test module user masih berfungsi normal

## 🧪 Testing Checklist

Setelah setiap module selesai direfactor:

- [ ] **Build Success**: `go build ./cmd/api` berhasil
- [ ] **No Import Errors**: Tidak ada error import
- [ ] **API Endpoints**: Test endpoint module masih berfungsi
- [ ] **Validation**: Validation rules masih bekerja
- [ ] **Handler Logic**: Handler logic masih berfungsi normal

## 📝 Final Verification

Setelah semua modules selesai:

- [ ] **Full Build**: `go build ./...` berhasil
- [ ] **Server Start**: Server bisa start tanpa error
- [ ] **All Endpoints**: Semua API endpoints masih berfungsi
- [ ] **Documentation Update**: Update dokumentasi jika diperlukan

## 🚨 Rollback Plan

Jika ada masalah:
1. Git revert ke commit sebelum refactoring
2. Atau restore individual files dari backup
3. Fix issues dan lanjutkan dari module yang bermasalah

## 📊 Expected Benefits

- **Reduced File Count**: 63 files → 45 files (28% reduction)
- **Faster Navigation**: Less file switching untuk developer
- **Cleaner Structure**: Logical grouping of related code
- **Easier Onboarding**: New developers less overwhelmed
- **Maintained Modularity**: Zero impact ke cross-module dependencies

---

**Status**: 🟡 Ready to Start
**Estimated Time**: ~2-3 hours untuk semua modules
**Risk Level**: 🟢 Low (no external dependencies affected)