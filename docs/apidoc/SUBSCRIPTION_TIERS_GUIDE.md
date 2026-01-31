# Subscription Tiers & Module Access Guide

## 📋 Overview

Sistem Huminor RBAC menggunakan subscription-based access control yang menentukan modul mana yang dapat diakses user berdasarkan subscription plan perusahaan.

## 🎯 Subscription Tiers

### **1. Basic Plan**
- **Price**: Rp 99,000/month
- **Max Users**: 25
- **Max Branches**: 3
- **Module Access**: Basic tier modules only

### **2. Professional Plan**
- **Price**: Rp 299,000/month
- **Max Users**: 100
- **Max Branches**: 10
- **Module Access**: Basic + Professional tier modules

### **3. Enterprise Plan**
- **Price**: Rp 599,000/month
- **Max Users**: Unlimited
- **Max Branches**: Unlimited
- **Module Access**: All modules (Basic + Professional + Enterprise)

### **4. Lifetime Plan** ⭐ **NEW**
- **Price**: One-time payment (varies)
- **Max Users**: Unlimited
- **Max Branches**: Unlimited
- **Module Access**: All modules (same as Enterprise)
- **Billing Cycle**: `lifetime`
- **Expiry**: Never expires (end_date: 2099-12-31)
- **Special Features**:
  - ✅ No recurring payments
  - ✅ Never expires
  - ✅ All current and future modules
  - ✅ Priority support

## 🔒 Module Access Logic

### **Active Subscription (Including Lifetime)**
```
User Login → Check Company Subscription → Get Plan Modules → Filter by Role → Return Modules
```

**Query Logic**:
```sql
SELECT DISTINCT m.*
FROM user_roles ur
JOIN role_modules rm ON ur.role_id = rm.role_id
JOIN modules m ON rm.module_id = m.id
JOIN plan_modules pm ON m.id = pm.module_id AND pm.is_included = true
JOIN subscriptions s ON pm.plan_id = s.plan_id
WHERE ur.user_id = ? 
  AND rm.can_read = true
  AND m.is_active = true
  AND s.company_id = ?
  AND s.status = 'active'
  AND (s.billing_cycle = 'lifetime' OR s.end_date >= CURRENT_DATE)
```

**Lifetime Subscription Handling**:
- ✅ `billing_cycle = 'lifetime'` bypasses expiry check
- ✅ `computed_status = 'lifetime'` in login response
- ✅ `days_remaining = null` (no expiration)
- ✅ Never falls back to basic modules

### **Expired Subscription**
```
User Login → Check Subscription (Expired) → Return Basic Modules Only
```

**Fallback Query Logic**:
```sql
SELECT DISTINCT m.*
FROM user_roles ur
JOIN role_modules rm ON ur.role_id = rm.role_id
JOIN modules m ON rm.module_id = m.id
WHERE ur.user_id = ? 
  AND rm.can_read = true
  AND m.is_active = true
  AND (m.subscription_tier = 'basic' OR m.subscription_tier IS NULL)
```

## 📊 Module Categories by Tier

### **Basic Tier Modules**
- **System & Security**: API Documentation, User Management, Role Management
- **Employee Self Service**: Basic Profile, Update Profile
- **Master Data**: Employee Data (basic fields)

### **Professional Tier Modules**
- **Attendance & Time**: Attendance System, Leave Management
- **Payroll & Compensation**: Basic Payroll, Salary Structure
- **Reporting & Analytics**: Basic Reports

### **Enterprise Tier Modules**
- **Module Management**: Full module control
- **Advanced Analytics**: Custom Reports, Dashboard
- **Performance Management**: KPI, Performance Review
- **Asset & Facility**: Asset Management, Facility Booking

## 🔄 Subscription Status Flow

### **Regular Subscriptions (Monthly/Yearly)**
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   ACTIVE        │───▶│   EXPIRED       │───▶│   RENEWED       │
│ All Plan Modules│    │ Basic Only      │    │ All Plan Modules│
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Full Access     │    │ Limited Access  │    │ Full Access     │
│ Based on Role   │    │ Basic Modules   │    │ Based on Role   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### **Lifetime Subscriptions** ⭐
```
┌─────────────────┐
│   LIFETIME      │ ◄─── Never expires
│ All Plan Modules│ ◄─── Always active
└─────────────────┘ ◄─── No renewal needed
         │
         ▼
┌─────────────────┐
│ Permanent Access│ ◄─── computed_status: 'lifetime'
│ Based on Role   │ ◄─── days_remaining: null
└─────────────────┘ ◄─── end_date: 2099-12-31
```

### **Computed Status Values**
- `active`: Regular subscription, not expiring soon
- `expiring_soon`: Expires within 7 days
- `expiring_today`: Expires today
- `expired`: Past end_date
- `lifetime`: ⭐ **NEW** - Never expires

## 🛡️ Security Implementation

### **1. Database Level**
- `subscriptions` table tracks company subscription status
- `plan_modules` table defines which modules are included in each plan
- `modules.subscription_tier` field categorizes modules by tier

### **2. Application Level**
- Auth repository checks subscription status on every login
- Expired subscriptions trigger fallback to basic modules only
- **Lifetime subscriptions bypass expiry checks**
- API endpoints validate module access based on subscription
- **Special handling for `billing_cycle = 'lifetime'`**

### **3. Frontend Level**
- Login response contains only accessible modules
- UI dynamically shows/hides features based on available modules
- Navigation menu filtered by subscription tier

## 📝 Testing Subscription Behavior

### **Test Active Subscription**
```bash
# Login with active subscription
curl -X POST "http://localhost:8081/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"user_identity": "800000001", "password": "password123"}' | jq '.data.user.modules'
```

**Expected**: All modules based on subscription plan + role assignment

### **Test Expired Subscription**
```sql
-- Expire subscription
UPDATE subscriptions SET end_date = CURRENT_DATE - INTERVAL '1 day' WHERE company_id = 1;
```

```bash
# Login with expired subscription
curl -X POST "http://localhost:8081/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"user_identity": "800000001", "password": "password123"}' | jq '.data.user.modules'
```

**Expected**: Only basic tier modules

### **Test No Subscription**
```sql
-- Remove subscription
DELETE FROM subscriptions WHERE company_id = 1;
```

**Expected**: Empty modules object `{}`

## 🔧 Troubleshooting

### **User Gets No Modules**
1. Check subscription status: `SELECT * FROM subscriptions WHERE company_id = ?`
2. Check subscription dates: `end_date >= CURRENT_DATE`
3. Check plan modules: `SELECT * FROM plan_modules WHERE plan_id = ?`
4. Check role assignments: `SELECT * FROM user_roles WHERE user_id = ?`

### **User Gets Wrong Modules**
1. Verify subscription plan: Basic/Pro/Enterprise
2. Check module tier: `SELECT subscription_tier FROM modules WHERE id = ?`
3. Verify plan includes module: `SELECT * FROM plan_modules WHERE plan_id = ? AND module_id = ?`

### **Subscription Expired But Still Getting Modules**
1. Check if modules have `subscription_tier = 'basic'` (this is expected)
2. Verify fallback logic is working correctly
3. Check if subscription end_date is properly set

## 💡 Best Practices

### **For Administrators**
- Monitor subscription expiry dates
- Set up alerts for upcoming renewals
- Regularly audit module access permissions
- Test subscription behavior in staging environment

### **For Developers**
- Always check subscription status in module access logic
- Implement graceful degradation for expired subscriptions
- Cache subscription status for performance
- Log subscription-related access attempts

### **For Business**
- Plan module tiers based on customer needs
- Provide clear upgrade paths between tiers
- Communicate subscription benefits clearly
- Monitor usage patterns by tier

---

**🎉 Summary**: Subscription system ensures proper access control while providing graceful degradation for expired subscriptions, maintaining basic functionality while encouraging renewals.
## 🧪 Testing Lifetime Subscription

### **Create Lifetime Plan**
```sql
INSERT INTO subscription_plans (name, display_name, description, price_monthly, price_yearly, max_users, max_branches, features, is_active)
VALUES (
  'lifetime',
  'Lifetime Plan',
  'Complete HR solution with lifetime access - no expiration',
  0.00,
  0.00,
  NULL,
  NULL,
  '{"lifetime_access": true, "no_expiration": true, "all_features": true}',
  true
);
```

### **Copy Modules from Enterprise Plan**
```sql
INSERT INTO plan_modules (plan_id, module_id, is_included)
SELECT 4, module_id, is_included
FROM plan_modules 
WHERE plan_id = 3 AND is_included = true;
```

### **Create Lifetime Subscription**
```sql
-- First, update billing_cycle constraint
ALTER TABLE subscriptions DROP CONSTRAINT subscriptions_billing_cycle_check;
ALTER TABLE subscriptions ADD CONSTRAINT subscriptions_billing_cycle_check 
CHECK (billing_cycle::text = ANY (ARRAY['monthly'::character varying, 'yearly'::character varying, 'lifetime'::character varying]::text[]));

-- Create lifetime subscription
INSERT INTO subscriptions (company_id, plan_id, status, billing_cycle, start_date, end_date, price, currency, payment_status, auto_renew)
VALUES (
  1, 
  4, 
  'active', 
  'lifetime', 
  '2026-01-31', 
  '2099-12-31',  -- Far future date
  0.00, 
  'IDR', 
  'paid', 
  false
);
```

### **Test Login Response**
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"user_identity": "800000001", "password": "password123"}' | jq '.data.user.subscription'
```

**Expected Response**:
```json
{
  "has_subscription": true,
  "subscription": {
    "computed_status": "lifetime",
    "days_remaining": null,
    "plan": {
      "name": "lifetime",
      "billing_cycle": "lifetime"
    },
    "end_date": "2099-12-31T00:00:00Z"
  }
}
```

## 🔧 Implementation Details

### **Database Changes**
1. **New Plan**: Added `lifetime` plan with all modules
2. **Constraint Update**: Added `'lifetime'` to billing_cycle constraint
3. **Far Future Date**: Set end_date to 2099-12-31 for lifetime subscriptions

### **Code Changes**
1. **Auth Repository**: Updated subscription queries to handle lifetime
2. **Computed Status**: Added `'lifetime'` status calculation
3. **Days Remaining**: Returns `null` for lifetime subscriptions
4. **Module Access**: Lifetime subscriptions bypass expiry checks

### **Query Logic Updates**
```sql
-- Before (only date check)
AND s.end_date >= CURRENT_DATE

-- After (lifetime OR date check)
AND (s.billing_cycle = 'lifetime' OR s.end_date >= CURRENT_DATE)
```

### **Status Computation**
```sql
CASE 
    WHEN s.billing_cycle = 'lifetime' THEN 'lifetime'
    WHEN s.end_date < CURRENT_DATE THEN 'expired'
    WHEN s.end_date = CURRENT_DATE THEN 'expiring_today'
    WHEN s.end_date <= CURRENT_DATE + INTERVAL '7 days' THEN 'expiring_soon'
    ELSE 'active'
END as computed_status
```

---

**Last Updated**: January 31, 2026  
**Version**: 2.0.0 (Added Lifetime Support)