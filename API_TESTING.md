# Demo Request Data - API Testing Examples

Contoh data dummy untuk demonstrasi dan testing semua endpoint API.  
**Copy-paste JSON di bawah ini untuk testing endpoint!**

---

## 📋 Table of Contents
1. [Authentication](#authentication)
2. [Users Management](#users-management)
3. [Achievements](#achievements)
4. [Students & Lecturers](#students--lecturers)
5. [Reports & Analytics](#reports--analytics)

---

## Authentication

### 1. POST /api/v1/auth/login
**Login sebagai Admin:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Login sebagai Mahasiswa:**
```json
{
  "username": "student001",
  "password": "student123"
}
```

**Login sebagai Dosen Wali:**
```json
{
  "username": "lecturer001",
  "password": "lecturer123"
}
```

---

### 2. POST /api/v1/auth/refresh
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDM1ODI0MDAsInVzZXJfaWQiOiJiNmE5NDIxOC1iZDk5LTRmY2YtOTJmOS0xZGQzM2YzNTU1NjQifQ.xxxxx"
}
```

---

### 3. POST /api/v1/auth/logout
```json
{}
```
*(No body required, hanya butuh Authorization header)*

---

### 4. GET /api/v1/auth/profile
*(No body required, hanya butuh Authorization header)*

---

## Users Management

### 5. GET /api/v1/users
**Query Parameters:**
```
?page=1&limit=10&search=admin
```

---

### 6. GET /api/v1/users/:id
**URL Example:**
```
/api/v1/users/b6a94218-bd99-4fcf-92f9-1dd33f355564
```

---

### 7. POST /api/v1/users
**Create Admin User (RECOMMENDED - pakai role_name):**
```json
{
  "username": "admin_baru",
  "email": "admin.baru@university.ac.id",
  "password": "Password123!",
  "full_name": "Admin Baru",
  "role_name": "Admin"
}
```

**Create Mahasiswa User (RECOMMENDED - auto detect dari student_id):**
```json
{
  "username": "mahasiswa_baru",
  "email": "mahasiswa.baru@university.ac.id",
  "password": "Password123!",
  "full_name": "Ahmad Rizki",
  "student_id": "205150200111999",
  "program_study": "Teknik Informatika",
  "academic_year": "2024"
}
```

**Create Mahasiswa User (alternatif - pakai role_name):**
```json
{
  "username": "mahasiswa_baru2",
  "email": "mahasiswa.baru2@university.ac.id",
  "password": "Password123!",
  "full_name": "Siti Aminah",
  "role_name": "Mahasiswa",
  "student_id": "205150200111998",
  "program_study": "Sistem Informasi",
  "academic_year": "2024"
}
```

**Create Dosen User (RECOMMENDED - auto detect dari lecturer_id):**
```json
{
  "username": "dosen_baru",
  "email": "dosen.baru@university.ac.id",
  "password": "Password123!",
  "full_name": "Dr. Budi Santoso, M.Kom",
  "lecturer_id": "198505152012121001",
  "department": "Teknik Informatika"
}
```

**Create Dosen User (alternatif - pakai role_name):**
```json
{
  "username": "dosen_baru2",
  "email": "dosen.baru2@university.ac.id",
  "password": "Password123!",
  "full_name": "Dr. Ani Wijaya, M.T",
  "role_name": "Dosen Wali",
  "lecturer_id": "199001012015041001",
  "department": "Sistem Informasi"
}
```

---

### 8. PUT /api/v1/users/:id
**Update User Profile:**
```json
{
  "email": "email.baru@university.ac.id",
  "full_name": "Nama Lengkap Baru",
  "is_active": true
}
```

**Nonaktifkan User:**
```json
{
  "is_active": false
}
```

**Update Email Only:**
```json
{
  "email": "email.updated@university.ac.id"
}
```

---

### 9. DELETE /api/v1/users/:id
**URL Example:**
```
DELETE /api/v1/users/b6a94218-bd99-4fcf-92f9-1dd33f355564
```
*(No body required)*

---

### 10. PUT /api/v1/users/:id/role
**Assign Role ke User:**
```json
{
  "role_id": "394ff7a7-2653-45e7-9f61-0deff113c345"
}
```

---

## Achievements

### 11. GET /api/v1/achievements
**Query Parameters:**
```
?page=1&limit=10&status=verified&achievement_type=competition
```

**Filter Options:**
- `status`: draft, submitted, verified, rejected
- `achievement_type`: competition, publication, organization, certification, academic, other
- `search`: search by title
- `student_id`: filter by student UUID

---

### 12. GET /api/v1/achievements/:id
**URL Example:**
```
/api/v1/achievements/693d6046d6115cb20e966b7c
```

---

### 13. POST /api/v1/achievements
**Create Competition Achievement:**
```json
{
  "achievement_type": "competition",
  "title": "Lomba Programming Nasional 2025",
  "description": "Juara 1 Kompetisi Pemrograman tingkat Nasional yang diselenggarakan oleh Kementerian Pendidikan",
  "data": {
    "competition_name": "Lomba Programming Nasional 2025",
    "organizer": "Kementerian Pendidikan dan Kebudayaan",
    "competition_level": "national",
    "achievement_rank": "1st place",
    "date": "2025-08-15",
    "location": "Jakarta Convention Center",
    "number_of_participants": 150,
    "certificate_url": "/uploads/certificate_programming.pdf",
    "documentation_url": "/uploads/photo_programming.jpg"
  }
}
```

**Create Publication Achievement:**
```json
{
  "achievement_type": "publication",
  "title": "Penelitian Machine Learning untuk Prediksi Penyakit",
  "description": "Publikasi jurnal internasional tentang penerapan machine learning dalam bidang kesehatan",
  "data": {
    "publication_type": "journal",
    "title": "Machine Learning Approach for Disease Prediction",
    "authors": ["Ahmad Rizki", "Dr. Budi Santoso"],
    "journal_name": "International Journal of Medical Informatics",
    "publication_level": "international",
    "issn": "1386-5056",
    "publication_date": "2025-06-01",
    "doi": "10.1016/j.ijmedinf.2025.104567",
    "url": "https://doi.org/10.1016/j.ijmedinf.2025.104567",
    "citation_count": 5
  }
}
```

**Create Organization Achievement:**
```json
{
  "achievement_type": "organization",
  "title": "Ketua Himpunan Mahasiswa Teknik Informatika",
  "description": "Menjabat sebagai Ketua HMTI periode 2024-2025",
  "data": {
    "organization_name": "Himpunan Mahasiswa Teknik Informatika",
    "position": "Ketua",
    "organization_level": "university",
    "start_date": "2024-01-01",
    "end_date": "2025-12-31",
    "description": "Memimpin organisasi dengan 200+ anggota aktif",
    "certificate_url": "/uploads/sk_ketua_hmti.pdf"
  }
}
```

**Create Certification Achievement:**
```json
{
  "achievement_type": "certification",
  "title": "AWS Certified Solutions Architect",
  "description": "Sertifikasi profesional dari Amazon Web Services",
  "data": {
    "certification_name": "AWS Certified Solutions Architect - Associate",
    "issuing_organization": "Amazon Web Services",
    "certification_level": "international",
    "issue_date": "2025-03-15",
    "expiry_date": "2028-03-15",
    "credential_id": "AWS-ASA-123456789",
    "credential_url": "https://www.credly.com/badges/xxxxx",
    "certificate_url": "/uploads/aws_certificate.pdf"
  }
}
```

**Create Academic Achievement:**
```json
{
  "achievement_type": "academic",
  "title": "Dean's List Semester Genap 2024/2025",
  "description": "Masuk Dean's List dengan IPK 3.95",
  "data": {
    "achievement_name": "Dean's List",
    "semester": "Genap 2024/2025",
    "gpa": 3.95,
    "achievement_date": "2025-07-01",
    "certificate_url": "/uploads/deans_list_certificate.pdf"
  }
}
```

**Create Other Achievement:**
```json
{
  "achievement_type": "other",
  "title": "Volunteer Kegiatan Sosial Kemasyarakatan",
  "description": "Volunteer dalam program pengabdian masyarakat desa terpencil",
  "data": {
    "activity_name": "Program Pengabdian Masyarakat 2025",
    "organizer": "Universitas Airlangga",
    "activity_level": "regional",
    "date": "2025-05-20",
    "duration_hours": 120,
    "role": "Koordinator Tim",
    "certificate_url": "/uploads/volunteer_certificate.pdf"
  }
}
```

---

### 14. PUT /api/v1/achievements/:id
**Update Achievement (hanya bisa edit yang status=draft):**
```json
{
  "title": "Lomba Programming Nasional 2025 - UPDATED",
  "description": "Deskripsi yang telah diperbarui",
  "data": {
    "competition_name": "Lomba Programming Nasional 2025",
    "organizer": "Kemendikbud RI",
    "competition_level": "national",
    "achievement_rank": "1st place",
    "date": "2025-08-15",
    "location": "Jakarta Convention Center, Indonesia",
    "number_of_participants": 200,
    "certificate_url": "/uploads/certificate_updated.pdf",
    "documentation_url": "/uploads/photo_updated.jpg"
  }
}
```

---

### 15. DELETE /api/v1/achievements/:id
**URL Example:**
```
DELETE /api/v1/achievements/693d6046d6115cb20e966b7c
```
*(No body required - soft delete)*

---

### 16. POST /api/v1/achievements/:id/submit
**Submit untuk Verifikasi:**
```json
{}
```
*(No body required - achievement status akan berubah dari draft ke submitted)*

**URL Example:**
```
POST /api/v1/achievements/693d6046d6115cb20e966b7c/submit
```

---

### 17. POST /api/v1/achievements/:id/verify
**Verify Achievement (Dosen Wali only):**
```json
{
  "comments": "Achievement telah diverifikasi. Sertifikat valid dan sesuai dengan kriteria."
}
```

**Verify tanpa comment (optional):**
```json
{}
```

**URL Example:**
```
POST /api/v1/achievements/693d6046d6115cb20e966b7c/verify
```

---

### 18. POST /api/v1/achievements/:id/reject
**Reject Achievement (Dosen Wali only):**
```json
{
  "reason": "Sertifikat tidak jelas dan dokumentasi kurang lengkap. Mohon upload ulang dengan kualitas yang lebih baik."
}
```

**URL Example:**
```
POST /api/v1/achievements/693d6046d6115cb20e966b7c/reject
```

---

### 19. GET /api/v1/achievements/:id/history
**URL Example:**
```
/api/v1/achievements/693d6046d6115cb20e966b7c/history
```
*(No body required)*

---

### 20. POST /api/v1/achievements/:id/attachments
**Upload File (menggunakan multipart/form-data):**

**cURL Example:**
```bash
curl -X POST http://localhost:3000/api/v1/achievements/693d6046d6115cb20e966b7c/attachments \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/certificate.pdf" \
  -F "description=Sertifikat Juara 1"
```

**Postman:**
- Method: POST
- Body type: form-data
- Key: `file` (type: File)
- Value: Select file
- Key: `description` (type: Text)
- Value: Deskripsi file

---

## Students & Lecturers

### 21. GET /api/v1/students
**Query Parameters:**
```
?page=1&limit=10&search=ahmad&program_study=Teknik Informatika
```

---

### 22. GET /api/v1/students/:id
**URL Example:**
```
/api/v1/students/88bc10ce-17c5-49ba-a2e6-fbc8ab83664f
```

---

### 23. GET /api/v1/students/:id/achievements
**URL Example:**
```
/api/v1/students/88bc10ce-17c5-49ba-a2e6-fbc8ab83664f/achievements?status=verified
```

**Query Parameters:**
- `status`: draft, submitted, verified, rejected
- `achievement_type`: competition, publication, etc.

---

### 24. PUT /api/v1/students/:id/advisor
**Assign Dosen Wali:**
```json
{
  "advisor_id": "d1359b6d-4b28-4319-b09c-66377fe34a4b"
}
```

**Ganti Dosen Wali:**
```json
{
  "advisor_id": "new-lecturer-uuid-here"
}
```

---

### 25. GET /api/v1/lecturers
**Query Parameters:**
```
?page=1&limit=10&search=budi
```

---

### 26. GET /api/v1/lecturers/:id/advisees
**URL Example:**
```
/api/v1/lecturers/d1359b6d-4b28-4319-b09c-66377fe34a4b/advisees
```
*(Menampilkan semua mahasiswa bimbingan dosen)*

---

## Reports & Analytics

### 27. GET /api/v1/reports/statistics
*(No body required - Admin atau Dosen Wali)*

**Response berisi:**
- Total students, lecturers
- Achievement counts by status
- Achievement counts by type
- Total points

---

### 28. GET /api/v1/reports/student/:id
**URL Example:**
```
/api/v1/reports/student/88bc10ce-17c5-49ba-a2e6-fbc8ab83664f
```

*(No body required)*

**Response berisi:**
- Student info
- Advisor info
- Achievement statistics
- Achievements by type
- Total points
- Recent achievements

---

### 29. GET /api/v1/reports/top-students 🆕
**Query Parameters:**
```
?limit=10
```

*(Menampilkan leaderboard top students berdasarkan jumlah prestasi terverifikasi)*

---

### 30. GET /api/v1/reports/statistics/period 🆕
**Query Parameters:**
```
?start_date=2024-01-01&end_date=2025-12-31
```

*(Statistik prestasi berdasarkan periode waktu)*

---

### 31. GET /api/v1/reports/statistics/competition-levels 🆕
*(No body required)*

*(Distribusi prestasi berdasarkan tingkat kompetisi: local, regional, national, international)*

---

## 📝 Notes untuk Testing

### ⚠️ PENTING - Format Request Body
- **Create User**: Fields FLAT (tidak nested)
  - ✅ BENAR: `{"username": "...", "student_id": "...", "program_study": "..."}`
  - ❌ SALAH: `{"username": "...", "student_data": {"student_id": "..."}}`
- **Update User**: Hanya support `full_name`, `email`, `is_active` (TIDAK ADA password)
- **Verify Achievement**: Field `comments` (bukan `verification_note`)
- **Reject Achievement**: Field `reason` (required)

### Authentication Header
Semua endpoint yang Protected memerlukan header:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Role IDs (untuk reference)
**PENTING: Gunakan role_id ini atau gunakan role_name sebagai alternatif**
```
Admin:       524c9756-5c56-491b-bb7d-c51b00748367  atau role_name: "Admin"
Mahasiswa:   394ff7a7-2653-45e7-9f61-0deff113c345  atau role_name: "Mahasiswa"
Dosen Wali:  7be08b8a-030e-4d2a-bef4-b90d81479b0a  atau role_name: "Dosen Wali"
```

### Create User - Required Fields by Role

**3 Cara Assign Role:**
1. **AUTO DETECT** (RECOMMENDED): Sistem otomatis detect role dari field
   - Jika ada `student_id` → Role Mahasiswa
   - Jika ada `lecturer_id` → Role Dosen Wali
2. **Pakai role_name**: `"role_name": "Admin"` atau `"role_name": "Mahasiswa"` atau `"role_name": "Dosen Wali"`
3. **Pakai role_id**: UUID role dari database

**Required Fields:**

**Admin:**
- username, email, password, full_name
- role_name: "Admin" (atau role_id)

**Mahasiswa:**
- username, email, password, full_name
- student_id, program_study, academic_year
- Optional: role_name: "Mahasiswa" (auto-detect dari student_id)

**Dosen Wali:**
- username, email, password, full_name
- lecturer_id, department
- Optional: role_name: "Dosen Wali" (auto-detect dari lecturer_id)

### Competition Levels
- `local` - Lokal/Kampus
- `regional` - Regional/Provinsi
- `national` - Nasional
- `international` - Internasional

### Achievement Types
- `competition` - Kompetisi/Perlombaan
- `publication` - Publikasi Karya Ilmiah
- `organization` - Organisasi
- `certification` - Sertifikasi
- `academic` - Akademik (Dean's List, Beasiswa, dll)
- `other` - Lainnya

### Achievement Status Flow
```
draft → submitted → verified/rejected
```

---

## 🚀 Quick Testing Script (PowerShell)

**Login dan Save Token:**
```powershell
$baseUrl = "http://localhost:3000/api/v1"
$login = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method Post -Body (@{username="admin"; password="admin123"} | ConvertTo-Json) -ContentType "application/json"
$token = $login.data.token
Write-Host "Token: $token"
```

**Test Create Achievement:**
```powershell
$headers = @{Authorization="Bearer $token"}
$body = @{
    achievement_type = "competition"
    title = "Test Lomba Programming"
    description = "Testing create achievement"
    data = @{
        competition_name = "Test Competition"
        organizer = "Test Organizer"
        competition_level = "national"
        achievement_rank = "1st place"
        date = "2025-08-15"
    }
} | ConvertTo-Json

Invoke-RestMethod -Uri "$baseUrl/achievements" -Method Post -Headers $headers -Body $body -ContentType "application/json"
```

---

**File ini siap digunakan untuk demo video! Copy-paste JSON sesuai endpoint yang dijelaskan.** 📹✨
