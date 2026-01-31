# 🔗 Module Mapping Guide

Panduan lengkap untuk mapping module ke user dalam sistem RBAC Huminor.

## 📋 Overview

Module mapping adalah proses menghubungkan module (fitur) dengan user melalui role. Flow lengkapnya:

```
Module → Plan Modules → Role Modules → User Roles → Login Response
```

## 🔄 Module Mapping Flow

### **1. Module harus included di Subscription Plan**
```sql
-- Check modules included in plan
SELECT pm.*, m.name as module_name 
FROM plan_modules pm 
JOIN modules m ON pm.module_id = m.id 
WHERE pm.plan_id = 3 AND pm.is_included = true;
```

### **2. Module di-assign ke Role**
```sql
-- Check modules assigned to role
SELECT rm.*, m.name as module_name 
FROM role_modules rm 
JOIN modules m ON rm.module_id = m.id 
WHERE rm.role_id = 13;
```

### **3. User memiliki Role tersebut**
```sql
-- Check user roles
SELECT ur.*, r.name as role_name 
FROM user_roles ur 
JOIN roles r ON ur.role_id = r.id 
WHERE ur.user_id = 16;
```

### **4. Company memiliki Subscription aktif**
```sql
-- Check active subscription
SELECT s.*, sp.name as plan_name 
FROM subscriptions s 
JOIN subscription_plans sp ON s.plan_id = sp.id 
WHERE s.company_id = 1 AND s.status = 'active' AND s.end_date >= CURRENT_DATE;
```

## 🛠️ API Endpoints untuk Module Mapping

### **Option 1: Add Modules (RECOMMENDED)**

**Endpoint:** `POST /api/v1/role-management/role/{roleId}/modules`

**Behavior:** Menambahkan module baru tanpa menghapus yang sudah ada

```bash
curl -X POST http://localhost:8081/api/v1/role-management/role/13/modules \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "modules": [
      {
        "module_id": 139,
        "can_read": true,
        "can_write": true,
        "can_delete": true
      },
      {
        "module_id": 140,
        "can_read": true,
        "can_write": false,
        "can_delete": false
      }
    ]
  }'
```

**✅ Advantages:**
- Module lama tetap ada
- Aman dari data loss
- Bisa update permissions module yang sudah ada

### **Option 2: Remove Specific Modules**

**Endpoint:** `DELETE /api/v1/role-management/role/{roleId}/modules`

**Behavior:** Menghapus module tertentu saja

```bash
curl -X DELETE http://localhost:8081/api/v1/role-management/role/13/modules \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "module_ids": [140, 141]
  }'
```

**✅ Advantages:**
- Hapus module tertentu saja
- Module lain tetap ada
- Selective removal

### **Option 3: Replace All Modules (USE WITH CAUTION)**

**Endpoint:** `PUT /api/v1/role-management/role/{roleId}/modules`

**Behavior:** Mengganti semua module (hapus semua, insert baru)

```bash
curl -X PUT http://localhost:8081/api/v1/role-management/role/13/modules \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "modules": [
      {
        "module_id": 139,
        "can_read": true,
        "can_write": true,
        "can_delete": true
      }
    ]
  }'
```

**⚠️ Warnings:**
- Module yang tidak disebutkan akan HILANG
- Bisa menyebabkan data loss
- Gunakan hanya jika yakin ingin replace semua

## 📝 Step-by-Step Guide

### **Scenario: Menambahkan Module Baru ke User yang Sudah Ada**

#### **Step 1: Login dan Get Token**
```bash
TOKEN=$(curl -s -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"user_identity": "800000001", "password": "password123"}' | jq -r '.data.access_token')
```

#### **Step 2: Check Current Role Modules**
```bash
curl -X GET http://localhost:8081/api/v1/roles/13/permissions \
  -H "Authorization: Bearer $TOKEN" | jq '.data.modules'
```

#### **Step 3: Add New Modules**
```bash
curl -X POST http://localhost:8081/api/v1/role-management/role/13/modules \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "modules": [
      {
        "module_id": 145,
        "can_read": true,
        "can_write": true,
        "can_delete": true
      },
      {
        "module_id": 146,
        "can_read": true,
        "can_write": true,
        "can_delete": true
      }
    ]
  }'
```

#### **Step 4: Verify User Access**
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"user_identity": "800000001", "password": "password123"}' | jq '.data.user.modules'
```

## 🔍 Troubleshooting

### **Problem: Module tidak muncul di login response**

#### **Check 1: Module included di subscription plan?**
```bash
curl -X GET http://localhost:8081/api/v1/admin/plan-modules/3 \
  -H "Authorization: Bearer $TOKEN"
```

#### **Check 2: Module assigned ke role?**
```bash
curl -X GET http://localhost:8081/api/v1/roles/13/permissions \
  -H "Authorization: Bearer $TOKEN"
```

#### **Check 3: User punya role tersebut?**
```bash
curl -X GET http://localhost:8081/api/v1/role-management/user/16/roles \
  -H "Authorization: Bearer $TOKEN"
```

#### **Check 4: Subscription masih aktif?**
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"user_identity": "800000001", "password": "password123"}' | jq '.data.user.subscription'
```

## 📊 Database Tables Reference

### **plan_modules**
```sql
plan_id | module_id | is_included
--------|-----------|------------
   3    |    139    |    true     -- ✅ Module 139 included in enterprise plan
   3    |    140    |    true     -- ✅ Module 140 included in enterprise plan
   1    |    139    |    false    -- ❌ Module 139 NOT included in basic plan
```

### **role_modules**
```sql
role_id | module_id | can_read | can_write | can_delete
--------|-----------|----------|-----------|------------
   13   |    139    |   true   |   true    |    true
   13   |    140    |   true   |   false   |    false
   13   |    145    |   true   |   true    |    true
```

### **user_roles**
```sql
user_id | role_id | company_id | branch_id | unit_id
--------|---------|------------|-----------|--------
   16   |   13    |     1      |     1     |   16
```

### **subscriptions**
```sql
company_id | plan_id | status | start_date | end_date
-----------|---------|--------|------------|----------
     1     |    3    | active | 2025-12-28 | 2026-02-28
```

## 🎯 Best Practices

### **✅ DO:**
- Gunakan `POST` untuk menambah module baru
- Gunakan `DELETE` untuk hapus module tertentu
- Check current modules sebelum modify
- Verify user access setelah changes
- Pastikan module included di subscription plan

### **❌ DON'T:**
- Jangan gunakan `PUT` kecuali yakin ingin replace semua
- Jangan assign module yang tidak included di plan
- Jangan lupa check subscription status
- Jangan modify role module tanpa backup

## 📚 Related Documentation

- [API_WORKFLOW_GUIDE.md](./API_WORKFLOW_GUIDE.md) - Complete API workflow
- [SUBSCRIPTION_TIERS_GUIDE.md](./SUBSCRIPTION_TIERS_GUIDE.md) - Subscription system guide
- [Postman Collection](../HUMINOR_RBAC_API_MODULE_BASED.postman_collection.json) - API testing collection

---

**Last Updated:** January 31, 2026  
**Version:** 1.0.0