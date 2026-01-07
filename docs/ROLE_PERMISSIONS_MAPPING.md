# Role Permissions Mapping - ERP RBAC System

## Overview

Sistem ERP RBAC menggunakan role-based access control dengan 4 level permission untuk setiap module:
- **Read (R)**: Dapat melihat/membaca data
- **Write (W)**: Dapat membuat/mengubah data
- **Delete (D)**: Dapat menghapus data
- **Approve (A)**: Dapat menyetujui request/transaksi

## Role Hierarchy & Permissions

### 1. SUPER ADMIN
**Deskripsi**: Mengelola seluruh sistem, konfigurasi, keamanan, dan akses tanpa batas

| Module | Read | Write | Delete | Approve |
|--------|:----:|:-----:|:------:|:-------:|
| **ALL MODULES** | ✅ | ✅ | ✅ | ✅ |

**Total Access**: Akses penuh ke semua 124 modules dalam sistem

---

### 2. HR ADMIN
**Deskripsi**: Menjalankan operasional HR harian dan mengelola data serta proses HR

| Modul | Read | Write | Delete | Approve |
|-------|:----:|:-----:|:------:|:-------:|
| Asset & Facility | ✅ | ✅ | ✅ | ❌ |
| Attendance & Time | ✅ | ✅ | ✅ | ✅ |
| Dashboard & Analytic | ✅ | ❌ | ❌ | ❌ |
| Disciplinary & Relations | ✅ | ✅ | ✅ | ✅ |
| Employee Self Service | ✅ | ✅ | ✅ | ❌ |
| Leave Management | ✅ | ✅ | ✅ | ✅ |
| Master Data | ✅ | ✅ | ✅ | ❌ |
| Offboarding & Exit | ✅ | ✅ | ✅ | ✅ |
| Payroll & Compensation | ✅ | ❌ | ❌ | ❌ |
| Performance Management | ✅ | ✅ | ✅ | ✅ |
| Recruitment | ✅ | ✅ | ✅ | ❌ |
| Reporting & Analytics | ✅ | ✅ | ❌ | ❌ |
| System & Security | ❌ | ❌ | ❌ | ❌ |
| Training & Development | ✅ | ✅ | ✅ | ❌ |

**Key Responsibilities**:
- Operasional HR harian dengan akses write ke sebagian besar module
- Approval authority untuk attendance, leave, disciplinary, performance
- Read-only access ke payroll dan dashboard
- No access ke system & security

---

### 3. HR MANAGER
**Deskripsi**: Mengawasi, menyetujui, dan menganalisis proses HR secara strategis

| Modul | Read | Write | Delete | Approve |
|-------|:----:|:-----:|:------:|:-------:|
| Asset & Facility | ✅ | ❌ | ❌ | ❌ |
| Attendance & Time | ✅ | ❌ | ❌ | ✅ |
| Dashboard & Analytic | ✅ | ❌ | ❌ | ❌ |
| Disciplinary & Relations | ✅ | ❌ | ❌ | ✅ |
| Employee Self Service | ❌ | ❌ | ❌ | ❌ |
| Leave Management | ✅ | ❌ | ❌ | ✅ |
| Master Data | ✅ | ❌ | ❌ | ❌ |
| Offboarding & Exit | ✅ | ❌ | ❌ | ✅ |
| Payroll & Compensation | ✅ | ❌ | ❌ | ❌ |
| Performance Management | ✅ | ✅ | ❌ | ✅ |
| Recruitment | ✅ | ❌ | ❌ | ✅ |
| Reporting & Analytics | ✅ | ❌ | ❌ | ❌ |
| System & Security | ❌ | ❌ | ❌ | ❌ |
| Training & Development | ✅ | ❌ | ❌ | ❌ |

**Key Responsibilities**:
- Strategic oversight dengan fokus pada approval authority
- Write access hanya untuk performance management
- Extensive read access untuk monitoring dan analysis
- No operational write access (delegated to HR Admin)

---

### 4. PAYROLL OFFICER
**Deskripsi**: Mengelola penggajian, kompensasi, pajak, dan kepatuhan payroll

| Modul | Read | Write | Delete | Approve |
|-------|:----:|:-----:|:------:|:-------:|
| Attendance & Time | ✅ | ❌ | ❌ | ❌ |
| Leave Management | ✅ | ❌ | ❌ | ❌ |
| Payroll & Compensation | ✅ | ✅ | ❌ | ✅ |
| Reporting & Analytics | ✅ | ✅ | ❌ | ❌ |
| **Lainnya** | ❌ | ❌ | ❌ | ❌ |

**Key Responsibilities**:
- Full operational control over payroll & compensation
- Read access ke attendance dan leave untuk payroll calculation
- Report generation capabilities
- Approval authority untuk payroll transactions

---

### 5. LINE MANAGER
**Deskripsi**: Mengelola dan menyetujui aktivitas HR untuk tim yang dipimpin

| Modul | Read | Write | Delete | Approve |
|-------|:----:|:-----:|:------:|:-------:|
| Attendance & Time | ✅ | ❌ | ❌ | ✅ |
| Leave Management | ✅ | ❌ | ❌ | ✅ |
| Performance Management | ✅ | ✅ | ❌ | ✅ |
| Training & Development | ✅ | ❌ | ❌ | ❌ |
| Offboarding & Exit | ✅ | ❌ | ❌ | ✅ |
| Dashboard & Analytic | ✅ | ❌ | ❌ | ❌ |
| **Lainnya** | ❌ | ❌ | ❌ | ❌ |

**Key Responsibilities**:
- Team management dengan approval authority untuk attendance dan leave
- Performance management untuk direct reports
- Offboarding approval untuk team members
- Dashboard access untuk team analytics

---

### 6. EMPLOYEE
**Deskripsi**: Mengakses layanan HR mandiri dan mengelola data pribadi

| Modul | Read | Write | Delete | Approve |
|-------|:----:|:-----:|:------:|:-------:|
| Employee Self Service | ✅ | ✅ | ❌ | ❌ |
| Attendance & Time | ✅ | ❌ | ❌ | ❌ |
| Leave Management | ✅ | ❌ | ❌ | ❌ |
| Performance Management | ✅ | ✅ | ❌ | ❌ |
| Training & Development | ✅ | ❌ | ❌ | ❌ |
| Payroll (Slip Gaji) | ✅ | ❌ | ❌ | ❌ |
| **Lainnya** | ❌ | ❌ | ❌ | ❌ |

**Key Responsibilities**:
- Self-service capabilities untuk personal data management
- Performance self-assessment dan goal setting
- View-only access untuk attendance, leave, dan payslip
- Training enrollment dan progress tracking

---

### 7. RECRUITER
**Deskripsi**: Mengelola proses rekrutmen dari perencanaan hingga onboarding

| Modul | Read | Write | Delete | Approve |
|-------|:----:|:-----:|:------:|:-------:|
| Recruitment | ✅ | ✅ | ❌ | ❌ |
| Onboarding | ✅ | ✅ | ❌ | ❌ |
| Reporting & Analytics | ✅ | ❌ | ❌ | ❌ |
| **Lainnya** | ❌ | ❌ | ❌ | ❌ |

**Key Responsibilities**:
- End-to-end recruitment process management
- Candidate onboarding coordination
- Recruitment analytics dan reporting
- No approval authority (requires HR Admin/Manager approval)

---

### 8. ASSET OFFICER
**Deskripsi**: Mengelola aset, inventaris, dan fasilitas perusahaan

| Modul | Read | Write | Delete | Approve |
|-------|:----:|:-----:|:------:|:-------:|
| Asset & Facility | ✅ | ✅ | ✅ | ✅ |
| Offboarding (Asset Return) | ✅ | ❌ | ❌ | ✅ |
| Reporting & Analytics | ✅ | ❌ | ❌ | ❌ |
| **Lainnya** | ❌ | ❌ | ❌ | ❌ |

**Key Responsibilities**:
- Complete asset lifecycle management
- Asset assignment dan maintenance approval
- Asset return processing during offboarding
- Asset utilization reporting

---

### 9. AUDITOR
**Deskripsi**: Mengakses laporan dan data HR secara read-only untuk audit

| Modul | Read | Write | Delete | Approve |
|-------|:----:|:-----:|:------:|:-------:|
| Reporting & Analytics | ✅ | ❌ | ❌ | ❌ |
| Payroll & Compensation | ✅ | ❌ | ❌ | ❌ |
| Attendance & Time | ✅ | ❌ | ❌ | ❌ |
| Disciplinary | ✅ | ❌ | ❌ | ❌ |
| Audit Log | ✅ | ❌ | ❌ | ❌ |

**Key Responsibilities**:
- Comprehensive read-only access untuk audit purposes
- Financial dan operational data review
- Compliance monitoring
- Audit trail analysis

---

### 10. IT ADMIN
**Deskripsi**: Mengelola keamanan, user, dan konfigurasi teknis sistem

| Modul | Read | Write | Delete | Approve |
|-------|:----:|:-----:|:------:|:-------:|
| System & Security | ✅ | ✅ | ✅ | ❌ |
| Master Data (technical) | ✅ | ❌ | ❌ | ❌ |
| Audit Log | ✅ | ❌ | ❌ | ❌ |
| **Lainnya** | ❌ | ❌ | ❌ | ❌ |

**Key Responsibilities**:
- System configuration dan security management
- User account management
- Technical master data maintenance
- System audit log monitoring

---

## Permission Matrix Summary

| Role | Total Modules | Read Access | Write Access | Delete Access | Approve Access |
|------|:-------------:|:-----------:|:------------:|:-------------:|:--------------:|
| SUPER_ADMIN | 124 | 124 | 124 | 124 | 124 |
| HR_ADMIN | 94 | 94 | 75 | 69 | 45 |
| HR_MANAGER | 78 | 78 | 7 | 0 | 42 |
| PAYROLL_OFFICER | 35 | 35 | 18 | 0 | 11 |
| LINE_MANAGER | 44 | 44 | 7 | 0 | 28 |
| EMPLOYEE | 44 | 44 | 20 | 0 | 0 |
| RECRUITER | 14 | 14 | 8 | 0 | 0 |
| ASSET_OFFICER | 25 | 25 | 7 | 7 | 8 |
| AUDITOR | 36 | 36 | 0 | 0 | 0 |
| IT_ADMIN | 8 | 8 | 3 | 3 | 0 |

## Implementation Notes

### Database Structure
- **Table**: `role_modules`
- **Columns**: `role_id`, `module_id`, `can_read`, `can_write`, `can_delete`, `can_approve`
- **Constraints**: Unique constraint pada (role_id, module_id)

### Subscription Integration
- Module access dibatasi berdasarkan subscription plan company
- Role permissions di-filter berdasarkan modules yang tersedia di subscription
- Enterprise plan: Akses ke semua modules sesuai role
- Basic plan: Akses terbatas ke modules basic tier saja

### API Integration
- Login response menampilkan modules sesuai role dan subscription
- Module access validation pada setiap API endpoint
- Real-time permission checking berdasarkan role_modules mapping

### Security Considerations
- Principle of least privilege diterapkan pada semua role
- Approval workflow memerlukan role dengan `can_approve = true`
- System & Security modules hanya accessible oleh IT_ADMIN dan SUPER_ADMIN
- Audit trail untuk semua permission changes

---

**Last Updated**: January 2026  
**Version**: 1.0  
**Total Roles**: 10  
**Total Modules**: 124  
**Total Permissions**: 4 (Read, Write, Delete, Approve)