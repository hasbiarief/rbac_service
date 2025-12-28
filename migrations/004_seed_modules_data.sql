-- Migration 004: Seed modules data
-- Complete module system with 89+ modules across 12 categories

-- Core HR / Master Data (Basic Tier)
INSERT INTO modules (category, name, url, icon, description, subscription_tier, is_active) VALUES
('Core HR / Master Data', 'Data Karyawan', '/core-hr/employees', 'Users', 'Kelola data karyawan', 'basic', true),
('Core HR / Master Data', 'Struktur Organisasi', '/core-hr/organization', 'Building', 'Kelola struktur organisasi perusahaan', 'basic', true),
('Core HR / Master Data', 'Jabatan & Posisi', '/core-hr/positions', 'Briefcase', 'Kelola jabatan dan posisi kerja', 'basic', true),
('Core HR / Master Data', 'Departemen', '/core-hr/departments', 'Building2', 'Kelola departemen perusahaan', 'basic', true),
('Core HR / Master Data', 'Status Karyawan', '/core-hr/employment-status', 'UserCheck', 'Kelola status kepegawaian', 'basic', true),
('Core HR / Master Data', 'Grade & Level', '/core-hr/grades', 'TrendingUp', 'Kelola grade dan level karyawan', 'basic', true),
('Core HR / Master Data', 'Lokasi Kerja', '/core-hr/locations', 'MapPin', 'Kelola lokasi kerja', 'basic', true),
('Core HR / Master Data', 'Shift Kerja', '/core-hr/shifts', 'Clock', 'Kelola shift kerja', 'basic', true),
('Core HR / Master Data', 'Hari Libur', '/core-hr/holidays', 'Calendar', 'Kelola hari libur dan cuti bersama', 'basic', true),
('Core HR / Master Data', 'Dokumen Karyawan', '/core-hr/documents', 'FileText', 'Kelola dokumen karyawan', 'basic', true),
('Core HR / Master Data', 'Kontak Darurat', '/core-hr/emergency-contacts', 'Phone', 'Kelola kontak darurat karyawan', 'basic', true),
('Core HR / Master Data', 'Riwayat Pendidikan', '/core-hr/education', 'GraduationCap', 'Kelola riwayat pendidikan', 'basic', true),
('Core HR / Master Data', 'Riwayat Pekerjaan', '/core-hr/work-history', 'Briefcase', 'Kelola riwayat pekerjaan', 'basic', true),
('Core HR / Master Data', 'Keluarga', '/core-hr/family', 'Users', 'Kelola data keluarga karyawan', 'basic', true),
('Core HR / Master Data', 'Keterampilan', '/core-hr/skills', 'Award', 'Kelola keterampilan karyawan', 'basic', true),
('Core HR / Master Data', 'Sertifikasi', '/core-hr/certifications', 'Award', 'Kelola sertifikasi karyawan', 'basic', true),
('Core HR / Master Data', 'Bank & Rekening', '/core-hr/bank-accounts', 'CreditCard', 'Kelola rekening bank karyawan', 'basic', true);

-- Employee Self Service (Basic Tier)
INSERT INTO modules (category, name, url, icon, description, subscription_tier, is_active) VALUES
('Employee Self Service', 'Update Profil', '/ess/profile', 'User', 'Update profil pribadi', 'basic', true),
('Employee Self Service', 'Pengajuan Cuti & Izin', '/ess/requests', 'FileText', 'Ajukan cuti dan izin', 'basic', true),
('Employee Self Service', 'Slip Gaji', '/ess/payslip', 'Receipt', 'Lihat slip gaji', 'basic', true),
('Employee Self Service', 'Dokumen Pribadi', '/ess/documents', 'Folder', 'Kelola dokumen pribadi', 'basic', true),
('Employee Self Service', 'Pengumuman', '/ess/announcements', 'Megaphone', 'Lihat pengumuman perusahaan', 'basic', true),
('Employee Self Service', 'Feedback & Saran', '/ess/feedback', 'MessageSquare', 'Berikan feedback dan saran', 'basic', true);

-- Recruitment (Pro Tier)
INSERT INTO modules (category, name, url, icon, description, subscription_tier, is_active) VALUES
('Recruitment', 'Manpower Planning', '/recruitment/planning', 'Target', 'Perencanaan kebutuhan tenaga kerja', 'pro', true),
('Recruitment', 'Job Vacancy Management', '/recruitment/vacancy', 'Briefcase', 'Kelola lowongan pekerjaan', 'pro', true),
('Recruitment', 'Applicant Tracking System', '/recruitment/ats', 'Users', 'Sistem pelacakan pelamar kerja', 'pro', true),
('Recruitment', 'Interview & Assessment', '/recruitment/interview', 'MessageSquare', 'Kelola proses wawancara dan penilaian', 'pro', true),
('Recruitment', 'Offering & Hiring', '/recruitment/hiring', 'UserCheck', 'Kelola penawaran kerja dan proses hiring', 'pro', true),
('Recruitment', 'Onboarding Karyawan Baru', '/recruitment/onboarding', 'UserPlus', 'Proses orientasi karyawan baru', 'pro', true);

-- Attendance & Time (Pro Tier)
INSERT INTO modules (category, name, url, icon, description, subscription_tier, is_active) VALUES
('Attendance & Time', 'Absensi Harian', '/attendance/presence', 'Clock', 'Kelola absensi harian karyawan', 'pro', true),
('Attendance & Time', 'Lembur', '/attendance/overtime', 'Clock', 'Kelola lembur karyawan', 'pro', true),
('Attendance & Time', 'Keterlambatan', '/attendance/late', 'AlertCircle', 'Kelola keterlambatan karyawan', 'pro', true),
('Attendance & Time', 'Shift Management', '/attendance/shift', 'RotateCcw', 'Kelola jadwal shift karyawan', 'pro', true),
('Attendance & Time', 'Rekap Absensi', '/attendance/recap', 'BarChart3', 'Rekap dan laporan absensi', 'pro', true),
('Attendance & Time', 'Kalender Libur Nasional', '/attendance/holiday', 'Calendar', 'Kelola kalender libur nasional', 'pro', true);

-- Leave Management (Pro Tier)
INSERT INTO modules (category, name, url, icon, description, subscription_tier, is_active) VALUES
('Leave Management', 'Jenis Cuti', '/leave/types', 'Calendar', 'Kelola jenis-jenis cuti', 'pro', true),
('Leave Management', 'Pengajuan Cuti & Izin', '/leave/request', 'FileText', 'Kelola pengajuan cuti dan izin karyawan', 'pro', true),
('Leave Management', 'Approval Workflow', '/leave/approval', 'CheckCircle', 'Kelola alur persetujuan cuti', 'pro', true),
('Leave Management', 'Saldo & Kuota Cuti', '/leave/balance', 'BarChart3', 'Kelola saldo dan kuota cuti karyawan', 'pro', true),
('Leave Management', 'Carry Forward & Expired Cuti', '/leave/carry-forward', 'ArrowRight', 'Kelola carry forward dan expired cuti', 'pro', true),
('Leave Management', 'Kalender Cuti Tim', '/leave/calendar', 'Calendar', 'Kalender cuti tim dan departemen', 'pro', true);

-- Performance Management (Pro Tier)
INSERT INTO modules (category, name, url, icon, description, subscription_tier, is_active) VALUES
('Performance Management', 'KPI & Target', '/performance/kpi', 'Target', 'Kelola KPI dan target karyawan', 'pro', true),
('Performance Management', 'Performance Review', '/performance/review', 'Star', 'Kelola review performa karyawan', 'pro', true),
('Performance Management', '360 Degree Feedback', '/performance/360-feedback', 'Users', 'Sistem feedback 360 derajat', 'pro', true),
('Performance Management', 'Goal Setting', '/performance/goals', 'Flag', 'Kelola goal setting karyawan', 'pro', true),
('Performance Management', 'Career Development', '/performance/career', 'TrendingUp', 'Kelola pengembangan karir', 'pro', true),
('Performance Management', 'Succession Planning', '/performance/succession', 'Users', 'Perencanaan suksesi jabatan', 'pro', true);

-- Training & Development (Pro Tier)
INSERT INTO modules (category, name, url, icon, description, subscription_tier, is_active) VALUES
('Training & Development', 'Training Program', '/training/programs', 'BookOpen', 'Kelola program pelatihan', 'pro', true),
('Training & Development', 'Training Schedule', '/training/schedule', 'Calendar', 'Kelola jadwal pelatihan', 'pro', true),
('Training & Development', 'Training Evaluation', '/training/evaluation', 'CheckSquare', 'Evaluasi hasil pelatihan', 'pro', true),
('Training & Development', 'Competency Management', '/training/competency', 'Award', 'Kelola kompetensi karyawan', 'pro', true),
('Training & Development', 'Learning Management', '/training/learning', 'GraduationCap', 'Sistem manajemen pembelajaran', 'pro', true),
('Training & Development', 'Training Budget', '/training/budget', 'DollarSign', 'Kelola budget pelatihan', 'pro', true);

-- Payroll & Compensation (Enterprise Tier)
INSERT INTO modules (category, name, url, icon, description, subscription_tier, is_active) VALUES
('Payroll & Compensation', 'Struktur Gaji', '/payroll/salary-structure', 'DollarSign', 'Kelola struktur gaji', 'enterprise', true),
('Payroll & Compensation', 'Komponen Gaji', '/payroll/components', 'List', 'Kelola komponen gaji', 'enterprise', true),
('Payroll & Compensation', 'Proses Payroll', '/payroll/process', 'Calculator', 'Proses penggajian bulanan', 'enterprise', true),
('Payroll & Compensation', 'Slip Gaji', '/payroll/payslip', 'Receipt', 'Generate slip gaji', 'enterprise', true),
('Payroll & Compensation', 'Pajak PPh 21', '/payroll/tax', 'FileText', 'Kelola pajak PPh 21', 'enterprise', true),
('Payroll & Compensation', 'BPJS Ketenagakerjaan', '/payroll/bpjs-tk', 'Shield', 'Kelola BPJS Ketenagakerjaan', 'enterprise', true),
('Payroll & Compensation', 'BPJS Kesehatan', '/payroll/bpjs-kes', 'Heart', 'Kelola BPJS Kesehatan', 'enterprise', true),
('Payroll & Compensation', 'Tunjangan', '/payroll/allowances', 'Plus', 'Kelola tunjangan karyawan', 'enterprise', true),
('Payroll & Compensation', 'Potongan', '/payroll/deductions', 'Minus', 'Kelola potongan gaji', 'enterprise', true),
('Payroll & Compensation', 'Bonus & Insentif', '/payroll/bonus', 'Gift', 'Kelola bonus dan insentif', 'enterprise', true);

-- Asset & Facility (Enterprise Tier)
INSERT INTO modules (category, name, url, icon, description, subscription_tier, is_active) VALUES
('Asset & Facility', 'Asset Management', '/asset/management', 'Package', 'Kelola aset perusahaan', 'enterprise', true),
('Asset & Facility', 'Asset Assignment', '/asset/assignment', 'UserCheck', 'Kelola penugasan aset ke karyawan', 'enterprise', true),
('Asset & Facility', 'Asset Maintenance', '/asset/maintenance', 'Tool', 'Kelola maintenance aset', 'enterprise', true),
('Asset & Facility', 'Facility Booking', '/facility/booking', 'Calendar', 'Booking fasilitas perusahaan', 'enterprise', true),
('Asset & Facility', 'Vehicle Management', '/asset/vehicle', 'Truck', 'Kelola kendaraan perusahaan', 'enterprise', true),
('Asset & Facility', 'Inventory Management', '/asset/inventory', 'Package', 'Kelola inventori perusahaan', 'enterprise', true);

-- Disciplinary & Relations (Enterprise Tier)
INSERT INTO modules (category, name, url, icon, description, subscription_tier, is_active) VALUES
('Disciplinary & Relations', 'Disciplinary Action', '/disciplinary/action', 'AlertTriangle', 'Kelola tindakan disipliner', 'enterprise', true),
('Disciplinary & Relations', 'Employee Relations', '/disciplinary/relations', 'Users', 'Kelola hubungan karyawan', 'enterprise', true),
('Disciplinary & Relations', 'Grievance Management', '/disciplinary/grievance', 'MessageSquare', 'Kelola pengaduan karyawan', 'enterprise', true),
('Disciplinary & Relations', 'Investigation', '/disciplinary/investigation', 'Search', 'Kelola investigasi kasus', 'enterprise', true);

-- Offboarding & Exit (Enterprise Tier)
INSERT INTO modules (category, name, url, icon, description, subscription_tier, is_active) VALUES
('Offboarding & Exit', 'Resignation Process', '/offboarding/resignation', 'UserMinus', 'Proses pengunduran diri', 'enterprise', true),
('Offboarding & Exit', 'Exit Interview', '/offboarding/exit-interview', 'MessageSquare', 'Wawancara keluar karyawan', 'enterprise', true),
('Offboarding & Exit', 'Asset Return', '/offboarding/asset-return', 'RotateCcw', 'Pengembalian aset perusahaan', 'enterprise', true),
('Offboarding & Exit', 'Final Settlement', '/offboarding/settlement', 'Calculator', 'Perhitungan settlement akhir', 'enterprise', true),
('Offboarding & Exit', 'Alumni Network', '/offboarding/alumni', 'Users', 'Jaringan alumni karyawan', 'enterprise', true),
('Offboarding & Exit', 'Knowledge Transfer', '/offboarding/knowledge-transfer', 'BookOpen', 'Transfer pengetahuan', 'enterprise', true);

-- Reporting & Analytics (Enterprise Tier)
INSERT INTO modules (category, name, url, icon, description, subscription_tier, is_active) VALUES
('Reporting & Analytics', 'HR Dashboard', '/reporting/dashboard', 'BarChart3', 'Dashboard HR analytics', 'enterprise', true),
('Reporting & Analytics', 'Employee Reports', '/reporting/employee', 'Users', 'Laporan data karyawan', 'enterprise', true),
('Reporting & Analytics', 'Attendance Reports', '/reporting/attendance', 'Clock', 'Laporan absensi', 'enterprise', true),
('Reporting & Analytics', 'Payroll Reports', '/reporting/payroll', 'DollarSign', 'Laporan payroll', 'enterprise', true),
('Reporting & Analytics', 'Performance Reports', '/reporting/performance', 'TrendingUp', 'Laporan performa', 'enterprise', true),
('Reporting & Analytics', 'Custom Reports', '/reporting/custom', 'FileText', 'Laporan kustom', 'enterprise', true);

-- System & Security (Basic Tier - Admin only)
INSERT INTO modules (category, name, url, icon, description, subscription_tier, is_active) VALUES
('System & Security', 'User Management', '/system/users', 'Users', 'Kelola pengguna sistem', 'basic', true),
('System & Security', 'Role Management', '/system/roles', 'Shield', 'Kelola role dan permission', 'basic', true),
('System & Security', 'Audit Log', '/system/audit', 'FileText', 'Log aktivitas sistem', 'basic', true),
('System & Security', 'System Settings', '/system/settings', 'Settings', 'Pengaturan sistem', 'basic', true),
('System & Security', 'API Management', '/system/api', 'Code', 'Kelola API access', 'basic', true),
('System & Security', 'Backup & Restore', '/system/backup', 'Database', 'Backup dan restore data', 'basic', true);