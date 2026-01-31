# API Workflow Guide - Complete Setup Process

## Overview

Panduan lengkap untuk menggunakan Huminor RBAC API dari awal sampai user dapat login dan melihat modul yang tersedia. Dokumen ini menjelaskan urutan API yang perlu dipanggil untuk setup complete system.

## 🏗️ Architecture Overview

```
Company → Subscription Plan → Available Modules
    ↓           ↓                    ↓
Branch → Unit → Role → User → Login Response
```

**Key Relationship:**
- Company must have active subscription
- Subscription plan determines available modules
- User gets modules through: Role → Plan Modules → Subscription
- Only modules included in company's subscription plan will appear in login response

### **Subscription System Architecture**

```
┌────────────────────────────────────────────────────────────┐
│                    SUBSCRIPTION SYSTEM                     │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  ┌─────────────────┐    ┌─────────────────┐                │
│  │ Subscription    │    │ Plan            │                │
│  │ Plans           │───▶│ Modules         │                │
│  │                 │    │ Mapping         │                │
│  │ • Basic         │    │                 │                │
│  │ • Professional  │    │ plan_modules    │                │
│  │ • Enterprise    │    │ table           │                │
│  └─────────────────┘    └─────────────────┘                │
│           │                       │                        │
│           │                       ▼                        │
│           │              ┌─────────────────┐               │
│           │              │ Available       │               │
│           │              │ Modules List    │               │
│           │              │ (Filtered)      │               │
│           │              └─────────────────┘               │
│           │                       │                        │
│           ▼                       │                        │
│  ┌─────────────────┐              │                        │
│  │ Company         │              │                        │
│  │ Subscription    │              │                        │
│  │                 │              │                        │
│  │ subscriptions   │              │                        │
│  │ table           │              │                        │
│  └─────────────────┘              │                        │
│           │                       │                        │
│           └───────────────────────┼────────────────────────┤
│                                   │                        │
└───────────────────────────────────┼────────────────────────┘
                                    │
                                    ▼
┌────────────────────────────────────────────────────────────┐
│                      RBAC SYSTEM                           │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  ┌─────────────────┐    ┌─────────────────┐                │
│  │ Role Module     │    │ User Role       │                │
│  │ Assignment      │    │ Assignment      │                │
│  │                 │    │                 │                │
│  │ role_modules    │    │ user_roles      │                │
│  │ table           │    │ table           │                │
│  │                 │    │                 │                │
│  │ (Only modules   │    │                 │                │
│  │ from            │    │                 │                │
│  │ subscription)   │    │                 │                │
│  └─────────────────┘    └─────────────────┘                │
│           │                       │                        │
│           └───────────────────────┼────────────────────────┤
│                                   │                        │
└───────────────────────────────────┼────────────────────────┘
                                    │
                                    ▼
                           ┌─────────────────┐
                           │ Login Response  │
                           │ Module List     │
                           │                 │
                           │ (Intersection   │
                           │ of Subscription │
                           │ & Role Access)  │
                           └─────────────────┘
```

### **Module Filtering Process**

1. **All System Modules** (150+ modules available)
   ↓
2. **Subscription Filter** (Plan determines available modules)
   - Basic Plan: Modules 1-20
   - Professional Plan: Modules 1-50  
   - Enterprise Plan: All modules
   ↓
3. **Role Assignment Filter** (User's roles determine accessible modules)
   ↓
4. **Permission Filter** (read/write/update/delete/approve)
   ↓
5. **Final Module List** (Shown in login response)

## 📋 Prerequisites

- Server running di `http://localhost:8081`
- Authentication token dari admin user
- Postman atau HTTP client untuk testing

## 🚀 Complete Workflow

### **Phase 1: System Foundation Setup**

#### **1.1 Login sebagai Admin**
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "user_identity": "800000001",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 16,
      "name": "Hasbi Due",
      "role_assignments": [...]
    }
  }
}
```

**📝 Action:** Save `access_token` untuk request selanjutnya.

---

### **Phase 2: Company & Subscription Setup**

#### **2.1 Create Company (Optional - jika belum ada)**
```http
POST /api/v1/companies
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "PT. Example Company",
  "code": "EXAMPLE"
}
```

#### **2.2 Check Available Subscription Plans**
```http
GET /api/v1/admin/subscription-plans
Authorization: Bearer {access_token}
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "basic",
      "display_name": "Basic Plan",
      "price_monthly": 99000.00,
      "max_users": 25,
      "max_branches": 3,
      "description": "Core features for small teams"
    },
    {
      "id": 2,
      "name": "pro", 
      "display_name": "Professional Plan",
      "price_monthly": 299000.00,
      "max_users": 100,
      "max_branches": 10,
      "description": "Advanced features for growing companies"
    },
    {
      "id": 3,
      "name": "enterprise",
      "display_name": "Enterprise Plan", 
      "price_monthly": 599000.00,
      "max_users": null,
      "max_branches": null,
      "description": "Complete access for large organizations"
    }
  ]
}
```

#### **2.2.1 Check Modules Included in Each Plan**
```http
GET /api/v1/admin/plan-modules/{plan_id}
Authorization: Bearer {access_token}
```

**Example for Basic Plan (plan_id=1):**
```json
{
  "success": true,
  "data": [
    {
      "module_id": 1,
      "module_name": "User Management",
      "module_url": "/users",
      "is_included": true,
      "subscription_tier": "basic"
    },
    {
      "module_id": 2,
      "module_name": "Role Management",
      "module_url": "/roles", 
      "is_included": true,
      "subscription_tier": "basic"
    },
    {
      "module_id": 50,
      "module_name": "Advanced Analytics",
      "module_url": "/analytics",
      "is_included": false,
      "subscription_tier": "enterprise"
    }
  ]
}
```

**📝 Action:** Choose appropriate plan based on required modules.

#### **2.3 Create Subscription for Company**

**Option A: Regular Subscription (Monthly/Yearly)**
```http
POST /api/v1/subscriptions
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "company_id": 1,
  "plan_id": 2,
  "billing_cycle": "monthly",
  "start_date": "2026-01-25",
  "end_date": "2026-02-25",
  "status": "active",
  "payment_status": "paid"
}
```

**Option B: Lifetime Subscription** ⭐ **NEW**
```http
POST /api/v1/subscriptions
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "company_id": 1,
  "plan_id": 4,
  "billing_cycle": "lifetime",
  "start_date": "2026-01-31",
  "end_date": "2099-12-31",
  "status": "active",
  "payment_status": "paid",
  "price": 0.00
}
```

**📝 Action:** Company sekarang memiliki subscription yang menentukan modules yang tersedia. Lifetime subscription tidak pernah expire.

#### **2.4 Create Branch**
```http
POST /api/v1/branches
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "company_id": 1,
  "name": "Cabang Jakarta",
  "code": "JKT"
}
```

**📝 Action:** Save `company_id`, `subscription_id`, dan `branch_id` dari response.

---

### **Phase 3: Unit Management**

#### **3.1 Create Unit**
```http
POST /api/v1/units
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "branch_id": 1,
  "name": "IT Department",
  "code": "IT",
  "description": "Information Technology Department"
}
```

**📝 Action:** Save `unit_id` dari response.

---

### **Phase 4: Module Management & Subscription Verification**

#### **4.1 Get Available Modules for Subscription Plan**
```http
GET /api/v1/admin/plan-modules/{plan_id}
Authorization: Bearer {access_token}
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "module_id": 1,
      "module_name": "User Management",
      "module_url": "/users",
      "is_included": true
    },
    {
      "module_id": 2,
      "module_name": "Role Management",
      "module_url": "/roles", 
      "is_included": true
    }
  ]
}
```

#### **4.2 Get All Available Modules**
```http
GET /api/v1/modules?limit=50
Authorization: Bearer {access_token}
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "User Management",
      "url": "/users",
      "description": "Manage system users",
      "subscription_tier": "basic"
    },
    {
      "id": 2,
      "name": "Role Management", 
      "url": "/roles",
      "description": "Manage user roles",
      "subscription_tier": "basic"
    },
    {
      "id": 50,
      "name": "Advanced Analytics",
      "url": "/analytics",
      "description": "Advanced reporting and analytics",
      "subscription_tier": "enterprise"
    }
  ]
}
```

#### **4.3 Add Module to Subscription Plan (Optional)**
```http
POST /api/v1/admin/plan-modules/{plan_id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "modules": [
    {
      "module_id": 50,
      "is_included": true
    }
  ]
}
```

**📝 Action:** Hanya modules yang included dalam subscription plan yang akan muncul di login response.

---

### **Phase 5: Role Management**

#### **5.1 Create Role**
```http
POST /api/v1/roles
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "Project Manager",
  "description": "Manages company projects and team members"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 15,
    "name": "Project Manager",
    "description": "Manages company projects and team members"
  }
}
```

**📝 Action:** Save `role_id` dari response.

#### **5.2 Assign Modules to Role (Only Subscription-Included Modules)**

**Option A: Replace All Modules (REPLACE)**
```http
PUT /api/v1/role-management/role/{role_id}/modules
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "modules": [
    {
      "module_id": 1,
      "can_read": true,
      "can_write": true,
      "can_delete": false
    },
    {
      "module_id": 2,
      "can_read": true,
      "can_write": false,
      "can_delete": false
    }
  ]
}
```
⚠️ **Warning:** Ini akan **mengganti semua** module. Module yang tidak disebutkan akan dihapus.

**Option B: Add New Modules (APPEND) - Recommended**
```http
POST /api/v1/role-management/role/{role_id}/modules
Authorization: Bearer {access_token}
Content-Type: application/json

{
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
}
```
✅ **Recommended:** Menambahkan module baru tanpa menghapus yang sudah ada.

**Option C: Remove Specific Modules**
```http
DELETE /api/v1/role-management/role/{role_id}/modules
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "module_ids": [140, 141]
}
```
✅ Menghapus module tertentu saja, module lain tetap ada.

**📝 Action:** Gunakan **POST** untuk menambah module baru, **DELETE** untuk hapus module tertentu, **PUT** hanya jika ingin replace semua.

**⚠️ Important:** Hanya assign modules yang sudah included dalam subscription plan company. Modules yang tidak included tidak akan muncul di login response meskipun di-assign ke role.

---

### **Phase 6: User Management**

#### **6.1 Create New User**
```http
POST /api/v1/users
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "user_identity": "100000001",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 17,
    "name": "John Doe",
    "email": "john.doe@example.com",
    "user_identity": "100000001"
  }
}
```

**📝 Action:** Save `user_id` dari response.

#### **6.2 Assign Role to User**
```http
POST /api/v1/role-management/assign-user-role
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "user_id": 17,
  "role_id": 15,
  "company_id": 1,
  "branch_id": 1,
  "unit_id": 1
}
```

**📝 Action:** User sekarang memiliki role dengan akses ke modules yang sudah di-assign.

---

### **Phase 7: Verification & Testing**

#### **7.1 Login as New User**
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "user_identity": "100000001",
  "password": "password123"
}
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "a3f485997fd448775128c5b9f5011ee3...",
    "user": {
      "id": 17,
      "name": "John Doe",
      "email": "john.doe@example.com",
      "user_identity": "100000001",
      "modules": {
        "User Management": [
          ["User Management", "/users", "Parent", "Manage system users"]
        ],
        "Role Management": [
          ["Role Management", "/roles", "Parent", "Manage user roles"]
        ],
        "Project Management": [
          ["Project Management", "/projects", "Parent", "Manage company projects"]
        ]
      },
      "role_assignments": [
        {
          "assignment_id": 35,
          "role_id": 15,
          "role_name": "Project Manager",
          "unit_id": 1,
          "unit_name": "IT Department",
          "branch_id": 1,
          "branch_name": "Cabang Jakarta",
          "company_id": 1,
          "company_name": "PT. Example Company"
        }
      ],
      "total_roles": 1,
      "subscription": {
        "has_subscription": true,
        "subscription": {
          "id": 1,
          "company_id": 1,
          "company_name": "PT. Example Company",
          "plan": {
            "id": 2,
            "name": "professional",
            "description": "Professional HR solution for growing companies",
            "price": 2990000,
            "billing_cycle": "yearly"
          },
          "limits": {
            "max_users": 100,
            "max_branches": 10
          },
          "status": "active",
          "computed_status": "active",
          "start_date": "2026-01-01T00:00:00Z",
          "end_date": "2027-01-01T00:00:00Z",
          "days_remaining": 335,
          "created_at": "2026-01-01T10:00:00Z",
          "updated_at": "2026-01-01T10:00:00Z"
        }
      }
    }
  }
}
```

**Lifetime Subscription Response Example** ⭐:
```json
{
  "success": true,
  "data": {
    "user": {
      "subscription": {
        "has_subscription": true,
        "subscription": {
          "id": 12,
          "company_id": 1,
          "company_name": "PT. Example Company",
          "plan": {
            "id": 4,
            "name": "lifetime",
            "description": "Complete HR solution with lifetime access - no expiration",
            "price": 0,
            "billing_cycle": "lifetime"
          },
          "limits": {
            "max_users": null,
            "max_branches": null
          },
          "status": "active",
          "computed_status": "lifetime",
          "start_date": "2026-01-31T00:00:00Z",
          "end_date": "2099-12-31T00:00:00Z",
          "days_remaining": null,
          "created_at": "2026-01-31T23:11:15Z",
          "updated_at": "2026-01-31T23:11:15Z"
        }
      }
    }
  }
}
```

#### **7.2.1 Understanding Subscription Information in Login Response**

The login response now includes comprehensive subscription information under the `subscription` field:

**Subscription Fields Explained:**
- `has_subscription`: Boolean indicating if company has active subscription
- `subscription.id`: Unique subscription identifier
- `subscription.company_id`: Company associated with subscription
- `subscription.company_name`: Company name
- `subscription.plan`: Subscription plan details
  - `id`: Plan identifier
  - `name`: Plan name (basic, professional, enterprise)
  - `description`: Plan description
  - `price`: Subscription price
  - `billing_cycle`: monthly or yearly
- `subscription.limits`: Plan limitations
  - `max_users`: Maximum users allowed (null = unlimited)
  - `max_branches`: Maximum branches allowed (null = unlimited)
- `subscription.status`: Subscription status (active, expired, cancelled, etc.)
- `subscription.computed_status`: Real-time status
  - `active`: Subscription is active
  - `expiring_soon`: Expires within 7 days
  - `expiring_today`: Expires today
  - `expired`: Already expired
  - `lifetime`: ⭐ **NEW** - Never expires
- `subscription.days_remaining`: Days until expiration (null for lifetime)
- `subscription.start_date`: Subscription start date
- `subscription.end_date`: Subscription end date

**Important Notes:**
- If `has_subscription: false`, user will only get basic tier modules
- Expired subscriptions automatically fall back to basic tier
- Only modules included in the subscription plan will appear in `modules` array

#### **7.2 Test User Permissions**
```http
GET /api/v1/users
Authorization: Bearer {new_user_access_token}
```

**Expected:** User dapat akses endpoint sesuai permission yang di-assign.

---

## 🔄 Complete API Sequence

### **Quick Setup Script (cURL)**

```bash
#!/bin/bash

# Variables
BASE_URL="http://localhost:8081/api/v1"
ADMIN_TOKEN=""

# 1. Login as admin
echo "1. Login as admin..."
ADMIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"user_identity": "800000001", "password": "password123"}')

ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | jq -r '.data.access_token')
echo "Admin token: $ADMIN_TOKEN"

# 2. Create company
echo "2. Creating company..."
COMPANY_RESPONSE=$(curl -s -X POST "$BASE_URL/companies" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "PT. Example Company", "code": "EXAMPLE"}')

COMPANY_ID=$(echo $COMPANY_RESPONSE | jq -r '.data.id')
echo "Company ID: $COMPANY_ID"

# 3. Create subscription for company
echo "3. Creating subscription..."
SUBSCRIPTION_RESPONSE=$(curl -s -X POST "$BASE_URL/subscriptions" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"company_id\": $COMPANY_ID, \"plan_id\": 2, \"billing_cycle\": \"monthly\", \"start_date\": \"2026-01-25\", \"end_date\": \"2026-02-25\", \"status\": \"active\", \"payment_status\": \"paid\"}")

SUBSCRIPTION_ID=$(echo $SUBSCRIPTION_RESPONSE | jq -r '.data.id')
echo "Subscription ID: $SUBSCRIPTION_ID"

# 4. Create branch
echo "4. Creating branch..."
BRANCH_RESPONSE=$(curl -s -X POST "$BASE_URL/branches" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"company_id\": $COMPANY_ID, \"name\": \"Cabang Jakarta\", \"code\": \"JKT\"}")

BRANCH_ID=$(echo $BRANCH_RESPONSE | jq -r '.data.id')
echo "Branch ID: $BRANCH_ID"

# 5. Create unit
echo "5. Creating unit..."
UNIT_RESPONSE=$(curl -s -X POST "$BASE_URL/units" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"branch_id\": $BRANCH_ID, \"name\": \"IT Department\", \"code\": \"IT\", \"description\": \"Information Technology Department\"}")

UNIT_ID=$(echo $UNIT_RESPONSE | jq -r '.data.id')
echo "Unit ID: $UNIT_ID"

# 6. Create role
echo "6. Creating role..."
ROLE_RESPONSE=$(curl -s -X POST "$BASE_URL/roles" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Project Manager", "description": "Manages projects"}')

ROLE_ID=$(echo $ROLE_RESPONSE | jq -r '.data.id')
echo "Role ID: $ROLE_ID"

# 7. Check available modules for subscription plan
echo "7. Checking available modules for plan..."
curl -s -X GET "$BASE_URL/admin/plan-modules/2" \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# 8. Assign modules to role (only subscription-included modules)
echo "8. Assigning modules to role..."
curl -s -X PUT "$BASE_URL/role-management/role/$ROLE_ID/modules" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"modules": [{"module_id": 1, "can_read": true, "can_write": true, "can_delete": false}, {"module_id": 2, "can_read": true, "can_write": false, "can_delete": false}]}'

# 9. Create user
echo "9. Creating user..."
USER_RESPONSE=$(curl -s -X POST "$BASE_URL/users" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"John Doe\", \"email\": \"john@example.com\", \"user_identity\": \"100000001\", \"password\": \"password123\"}")

USER_ID=$(echo $USER_RESPONSE | jq -r '.data.id')
echo "User ID: $USER_ID"

# 10. Assign role to user
echo "10. Assigning role to user..."
curl -s -X POST "$BASE_URL/role-management/assign-user-role" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\": $USER_ID, \"role_id\": $ROLE_ID, \"company_id\": $COMPANY_ID, \"branch_id\": $BRANCH_ID, \"unit_id\": $UNIT_ID}"

# 11. Test login as new user
echo "11. Testing login as new user..."
USER_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"user_identity": "100000001", "password": "password123"}')

echo "User login response:"
echo $USER_LOGIN_RESPONSE | jq '.'

echo "Setup complete!"
```

---

## 🔧 **How to Map New Modules to Existing Users**

### **Scenario:** Anda sudah memiliki user dan role, tapi ingin menambahkan module baru

#### **Step 1: Check Current Role Modules**
```bash
# Get current modules for role
curl -X GET "$BASE_URL/roles/{role_id}/permissions" \
  -H "Authorization: Bearer $TOKEN"
```

#### **Step 2: Add New Modules (Recommended)**
```bash
# Add modules without removing existing ones
curl -X POST "$BASE_URL/role-management/role/{role_id}/modules" \
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

#### **Step 3: Remove Specific Modules (If Needed)**
```bash
# Remove specific modules only
curl -X DELETE "$BASE_URL/role-management/role/{role_id}/modules" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "module_ids": [140, 141]
  }'
```

#### **Step 4: Verify User Access**
```bash
# Login as user to check new modules
curl -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"user_identity": "800000001", "password": "password123"}' | jq '.data.user.modules'
```

### **⚠️ Important Notes:**

1. **Module harus included di subscription plan** - Check dengan:
   ```bash
   curl -X GET "$BASE_URL/admin/plan-modules/{plan_id}" \
     -H "Authorization: Bearer $TOKEN"
   ```

2. **Gunakan POST untuk ADD, bukan PUT** - PUT akan replace semua module

3. **Module mapping flow:**
   ```
   Module → Plan Modules (is_included=true) → Role Modules → User Roles → Login Response
   ```

---

## 📊 Complete Data Flow Diagram

### **Primary Flow: Company Setup & Subscription**
```
┌─────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Company   │───▶│ Subscription    │───▶│ Subscription    │
│             │    │ Plan Selection  │    │ Activation      │
└─────────────┘    └─────────────────┘    └─────────────────┘
       │                     │                       │
       │                     ▼                       │
       │            ┌─────────────────┐               │
       │            │ Available       │               │
       │            │ Modules         │◀──────────────┘
       │            │ (Plan-Based)    │
       │            └─────────────────┘
       │                     │
       ▼                     │
┌─────────────┐              │
│   Branch    │              │
└─────────────┘              │
       │                     │
       ▼                     │
┌─────────────┐              │
│    Unit     │              │
└─────────────┘              │
       │                     │
       ▼                     ▼
┌─────────────┐    ┌─────────────────┐
│    Role     │◀───│ Module          │
│             │    │ Assignment      │
└─────────────┘    │ (Filtered by    │
       │           │ Subscription)   │
       │           └─────────────────┘
       ▼
┌─────────────┐
│    User     │
│ Assignment  │
└─────────────┘
       │
       ▼
┌─────────────────┐
│ Login Response  │
│ with Modules    │
│ (Subscription   │
│ + Role Based)   │
└─────────────────┘
```

### **Subscription System Detail**
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Subscription    │    │ Plan            │    │ Plan            │
│ Plans           │───▶│ Modules         │───▶│ Available       │
│ (Basic/Pro/Ent) │    │ Mapping         │    │ Modules List    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                       │
                                ▼                       │
                       ┌─────────────────┐               │
                       │ Company         │               │
                       │ Subscription    │               │
                       │ (Active Plan)   │               │
                       └─────────────────┘               │
                                │                       │
                                ▼                       │
                       ┌─────────────────┐               │
                       │ Role Module     │◀──────────────┘
                       │ Assignment      │
                       │ (Only from      │
                       │ Subscribed      │
                       │ Modules)        │
                       └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │ User Login      │
                       │ Response        │
                       │ (Final Module   │
                       │ List)           │
                       └─────────────────┘
```

### **Module Visibility Logic**
```
All System Modules (1-150+)
           │
           ▼
    ┌─────────────────┐
    │ Subscription    │ ──── Filter 1: Plan Inclusion
    │ Plan Filter     │      (Basic: 1-20, Pro: 1-50, Enterprise: All)
    └─────────────────┘
           │
           ▼
    Available Modules for Company
           │
           ▼
    ┌─────────────────┐
    │ Role Module     │ ──── Filter 2: Role Assignment  
    │ Assignment      │      (User gets only assigned modules)
    └─────────────────┘
           │
           ▼
    ┌─────────────────┐
    │ Permission      │ ──── Filter 3: Permission Level
    │ Level Filter    │      (read/write/update/delete/approve)
    └─────────────────┘
           │
           ▼
    Final Module List in Login Response
```

---

## 🎯 Key Points

### **1. Hierarchy Importance**
- Company → Branch → Unit → Role → User
- Setiap level harus dibuat berurutan
- User assignment harus sesuai hierarchy

### **2. Subscription System (CRITICAL)**
- **Company MUST have active subscription**
- **Subscription plan determines available modules**
- **Module visibility flow**: All Modules → Plan Filter → Role Assignment → User Access
- **Three subscription tiers**:
  - **Basic Plan**: Modules 1-20 (Core features)
  - **Professional Plan**: Modules 1-50 (Advanced features)  
  - **Enterprise Plan**: All Modules (Complete access)

### **3. Module Assignment Logic**
- Modules di-assign ke Role, bukan langsung ke User
- **ONLY modules included in company's subscription plan can be assigned to roles**
- Permission level: read, write, update, delete, approve
- User mendapat akses module melalui: Subscription → Role → User

### **4. Login Response Structure**
- `access_token`: JWT token for API authentication
- `refresh_token`: Token for refreshing access token
- `user.modules`: Grouped modules yang bisa diakses user
- **Modules shown = Subscription Plan ∩ Role Assignment ∩ User Permission**
- `user.role_assignments`: Detail role dan unit assignment
- `user.total_roles`: Jumlah role yang dimiliki user
- `user.subscription`: **NEW** - Comprehensive subscription information
  - `has_subscription`: Boolean indicating active subscription
  - `subscription.plan`: Plan details (name, price, billing_cycle)
  - `subscription.limits`: Plan limitations (max_users, max_branches)
  - `subscription.computed_status`: Real-time status (active, expiring_soon, expired)
  - `subscription.days_remaining`: Days until expiration
  - `subscription.start_date` & `subscription.end_date`: Subscription period

### **5. Permission Checking**
- Setiap API call akan check permission berdasarkan role
- Module access ditentukan oleh: subscription_plan → role_modules → user_roles
- Unit-level isolation untuk data access

### **6. Subscription Database Tables**
```sql
-- Core subscription tables
subscription_plans (id, name, price_monthly, max_users, max_branches)
plan_modules (plan_id, module_id, is_included)
subscriptions (company_id, plan_id, status, start_date, end_date)

-- Module visibility query logic:
SELECT m.* FROM modules m
JOIN plan_modules pm ON m.id = pm.module_id  
JOIN subscriptions s ON pm.plan_id = s.plan_id
JOIN role_modules rm ON m.id = rm.module_id
JOIN user_roles ur ON rm.role_id = ur.role_id
WHERE s.company_id = ? AND s.status = 'active' 
  AND pm.is_included = true AND ur.user_id = ?
```

---

## 🔧 Troubleshooting

### **Common Issues:**

#### **1. User tidak dapat login**
- ✅ Check `is_active = true`
- ✅ Verify password
- ✅ Check user_identity format

#### **2. Module tidak muncul di login response (MOST COMMON)**
- ✅ **STEP 1: Check company has active subscription**
  ```sql
  SELECT * FROM subscriptions 
  WHERE company_id = ? AND status = 'active' 
  AND start_date <= NOW() AND end_date >= NOW();
  ```
- ✅ **STEP 2: Verify module is included in subscription plan**
  ```sql
  SELECT pm.*, m.name as module_name 
  FROM plan_modules pm
  JOIN modules m ON pm.module_id = m.id
  WHERE pm.plan_id = ? AND pm.is_included = true;
  ```
- ✅ **STEP 3: Check role assignment to user**
  ```sql
  SELECT ur.*, r.name as role_name 
  FROM user_roles ur
  JOIN roles r ON ur.role_id = r.id
  WHERE ur.user_id = ?;
  ```
- ✅ **STEP 4: Check module assignment to role**
  ```sql
  SELECT rm.*, m.name as module_name
  FROM role_modules rm
  JOIN modules m ON rm.module_id = m.id
  WHERE rm.role_id = ?;
  ```
- ✅ **STEP 5: Verify permission settings**

#### **3. Permission denied saat akses API**
- ✅ Check role memiliki module access
- ✅ Verify permission level (read/write/etc)
- ✅ Check unit-level access
- ✅ **Verify subscription includes the module**

#### **4. Empty modules di login response**
- ✅ **Company must have active subscription** (PRIMARY CHECK)
- ✅ **Subscription plan must include modules** (SECONDARY CHECK)
- ✅ User harus memiliki minimal 1 role
- ✅ Role harus memiliki minimal 1 module
- ✅ Module harus active

#### **5. Subscription-related issues**
- ✅ **Check subscription status**:
  ```sql
  SELECT s.*, sp.name as plan_name, sp.display_name,
    CASE 
      WHEN s.end_date >= CURRENT_DATE AND s.status = 'active' THEN 'ACTIVE'
      WHEN s.end_date < CURRENT_DATE THEN 'EXPIRED'
      WHEN s.status != 'active' THEN 'INACTIVE'
      ELSE 'UNKNOWN'
    END as computed_status
  FROM subscriptions s
  JOIN subscription_plans sp ON s.plan_id = sp.id
  WHERE s.company_id = ?;
  ```
- ✅ **Verify plan-module mapping**:
  ```sql
  SELECT pm.*, m.name as module_name, m.url
  FROM plan_modules pm
  JOIN modules m ON pm.module_id = m.id
  WHERE pm.plan_id = ? AND pm.is_included = true
  ORDER BY m.id;
  ```
- ✅ **Check subscription dates**:
  ```sql
  SELECT *, 
    CASE 
      WHEN start_date > NOW() THEN 'NOT_STARTED'
      WHEN end_date < NOW() THEN 'EXPIRED'
      WHEN status != 'active' THEN 'INACTIVE'
      ELSE 'ACTIVE'
    END as subscription_status
  FROM subscriptions WHERE company_id = ?;
  ```

#### **6. Subscription Expiry Behavior**
- ✅ **Active Subscription**: User gets all modules included in their subscription plan
- ✅ **Expired Subscription**: User only gets modules with `subscription_tier = 'basic'` or `subscription_tier IS NULL`
- ✅ **No Subscription**: User gets no modules (empty response)

**Example Behavior**:
```
Enterprise Plan (Active) → All modules (basic + pro + enterprise)
Enterprise Plan (Expired) → Only basic tier modules
No Subscription → No modules
```

#### **6. Module assignment issues**
- ✅ **Cannot assign module to role if not in subscription**:
  ```sql
  -- This query should return the module for successful assignment
  SELECT m.* FROM modules m
  JOIN plan_modules pm ON m.id = pm.module_id
  JOIN subscriptions s ON pm.plan_id = s.plan_id
  WHERE m.id = ? AND s.company_id = ? 
    AND s.status = 'active' AND pm.is_included = true;
  ```
- ✅ **Check role-module assignment**:
  ```sql
  SELECT rm.*, m.name, r.name as role_name
  FROM role_modules rm
  JOIN modules m ON rm.module_id = m.id
  JOIN roles r ON rm.role_id = r.id
  WHERE rm.role_id = ?;
  ```

### **Debug Commands for Subscription Issues**

#### **Complete Module Visibility Debug**
```sql
-- 1. Check user's company subscription
SELECT 
  c.name as company_name,
  s.status as subscription_status,
  sp.name as plan_name,
  sp.display_name as plan_display_name,
  s.start_date,
  s.end_date,
  CASE 
    WHEN s.start_date > NOW() THEN 'NOT_STARTED'
    WHEN s.end_date < NOW() THEN 'EXPIRED'  
    WHEN s.status != 'active' THEN 'INACTIVE'
    ELSE 'ACTIVE'
  END as computed_status
FROM users u
JOIN companies c ON u.company_id = c.id
JOIN subscriptions s ON c.id = s.company_id
JOIN subscription_plans sp ON s.plan_id = sp.id
WHERE u.id = ?;

-- 2. Check available modules for user's subscription plan
SELECT 
  m.id,
  m.name,
  m.url,
  pm.is_included,
  sp.name as plan_name
FROM users u
JOIN companies c ON u.company_id = c.id
JOIN subscriptions s ON c.id = s.company_id
JOIN subscription_plans sp ON s.plan_id = sp.id
JOIN plan_modules pm ON sp.id = pm.plan_id
JOIN modules m ON pm.module_id = m.id
WHERE u.id = ? AND pm.is_included = true
ORDER BY m.id;

-- 3. Check user's role assignments and module access
SELECT 
  ur.user_id,
  r.name as role_name,
  m.name as module_name,
  rm.can_read,
  rm.can_write,
  rm.can_update,
  rm.can_delete,
  rm.can_approve,
  pm.is_included as subscription_includes_module
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
JOIN role_modules rm ON r.id = rm.role_id
JOIN modules m ON rm.module_id = m.id
JOIN users u ON ur.user_id = u.id
JOIN subscriptions s ON u.company_id = s.company_id
JOIN plan_modules pm ON s.plan_id = pm.plan_id AND m.id = pm.module_id
WHERE ur.user_id = ? AND s.status = 'active' AND pm.is_included = true
ORDER BY m.id;
```

---

## 📚 Related Documentation

- [API Overview](API_OVERVIEW.md) - Complete API documentation
- [Authentication Integration Guide](../integration/AUTH_INTEGRATION_QUICK_START.md) - Integration guide
- [Security Best Practices](../integration/SECURITY_BEST_PRACTICES.md) - Security guidelines
- [Postman Collection](../HUMINOR_RBAC_API_MODULE_BASED.postman_collection.json) - Testing collection

---

**🎉 Congratulations!** Anda sekarang memiliki complete workflow untuk setup Huminor RBAC system dari awal sampai user dapat login dan melihat modules yang tersedia.