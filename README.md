# Student Achievement System

Sistem manajemen pencatatan dan verifikasi prestasi mahasiswa berbasis REST API menggunakan Go (Fiber Framework) dengan arsitektur Clean Architecture.

## 📋 Daftar Isi

- [Fitur Utama](#fitur-utama)
- [Teknologi](#teknologi)
- [Arsitektur](#arsitektur)
- [Instalasi](#instalasi)
- [Konfigurasi](#konfigurasi)
- [Menjalankan Aplikasi](#menjalankan-aplikasi)
- [Testing](#testing)
- [Dokumentasi API](#dokumentasi-api)
- [Role & Permission](#role--permission)
- [Struktur Database](#struktur-database)
- [Kontributor](#kontributor)

---

## 🎯 Fitur Utama

### Authentication & Authorization
- ✅ JWT-based authentication dengan refresh token
- ✅ Role-Based Access Control (RBAC)
- ✅ Permission-based authorization
- ✅ Secure password hashing (bcrypt)

### User Management
- ✅ CRUD pengguna (Admin, Mahasiswa, Dosen Wali)
- ✅ Auto-assign role berdasarkan tipe user
- ✅ Soft delete untuk data integrity
- ✅ Profile management

### Achievement Management
- ✅ CRUD prestasi mahasiswa
- ✅ Berbagai tipe prestasi (kompetisi, publikasi, research grant, dll)
- ✅ Flexible data structure dengan MongoDB
- ✅ Status tracking (draft, pending, verified, rejected)

### Verification Workflow
- ✅ Submit prestasi untuk verifikasi
- ✅ Approve/reject oleh dosen wali
- ✅ Comments dan feedback system
- ✅ Revision support untuk prestasi yang ditolak

### Reporting & Analytics
- ✅ Dashboard statistik
- ✅ Laporan per mahasiswa
- ✅ Filter berdasarkan status, level, tipe prestasi
- ✅ Timeline achievements

---

## 🛠 Teknologi

### Backend
- **Go 1.21+** - Programming language
- **Fiber v2** - Web framework
- **GORM** - ORM untuk PostgreSQL
- **MongoDB Driver** - NoSQL database driver
- **JWT-Go** - JSON Web Token
- **Bcrypt** - Password hashing

### Database
- **PostgreSQL** - Relational data (users, roles, permissions)
- **MongoDB** - Document storage (achievements, flexible data)

### Tools
- **Air** - Live reload (development)
- **Docker** - Containerization (optional)

---

## 🏗 Arsitektur

### Clean Architecture (2-Layer)

```
┌─────────────────────────────────────────┐
│           HTTP Handler Layer            │
│  (Routes + Middleware + Controllers)    │
└─────────────────┬───────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│          Service Layer                  │
│    (Business Logic + Validation)        │
└─────────────────┬───────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│        Repository Layer                 │
│     (Database Access + Queries)         │
└─────────────────┬───────────────────────┘
                  │
         ┌────────┴────────┐
         ▼                 ▼
   ┌──────────┐      ┌──────────┐
   │PostgreSQL│      │ MongoDB  │
   └──────────┘      └──────────┘
```

### Struktur Folder

```
Project-UAS-Advanced-Backend/
├── config/              # Konfigurasi aplikasi
│   └── config.go
├── database/            # Database initialization & seeding
│   ├── postgres.go
│   ├── mongodb.go
│   └── seed.go
├── middleware/          # Authentication & Authorization
│   ├── auth.go
│   └── permission.go
├── models/              # Data models
│   ├── user.go
│   ├── student.go
│   └── achievement.go
├── repository/          # Database layer
│   ├── user_repository.go
│   ├── student_repository.go
│   ├── lecturer_repository.go
│   └── achievement_repository.go
├── routes/              # Route definitions
│   └── routes.go
├── service/             # Business logic layer
│   ├── auth_service.go
│   ├── user_service.go
│   ├── achievement_service.go
│   ├── verification_service.go
│   └── report_service.go
├── utils/               # Helper functions
│   ├── response.go
│   ├── jwt.go
│   └── password.go
├── main.go              # Application entry point
├── .env                 # Environment variables
└── go.mod               # Go dependencies
```

---

## 📦 Instalasi

### Prerequisites

1. **Go 1.21 atau lebih tinggi**
   ```bash
   go version
   ```

2. **PostgreSQL 14+**
   ```bash
   psql --version
   ```

3. **MongoDB 6.0+**
   ```bash
   mongod --version
   ```

### Clone Repository

```bash
git clone https://github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend.git
cd Project-UAS-Advanced-Backend
```

### Install Dependencies

```bash
go mod download
```

---

## ⚙️ Konfigurasi

### 1. Setup Database

#### PostgreSQL
```bash
# Login ke PostgreSQL
psql -U postgres

# Buat database
CREATE DATABASE achievement_db;

# Keluar
\q
```

#### MongoDB
```bash
# MongoDB akan otomatis membuat database saat pertama kali digunakan
# Pastikan MongoDB service berjalan
```

### 2. Environment Variables

Buat file `.env` di root directory:

```env
# Application
APP_NAME=Student Achievement System
APP_ENV=development
APP_PORT=3000
API_VERSION=v1

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=admin
POSTGRES_DB=achievement_db

# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=achievement_db

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_REFRESH_SECRET=your-super-secret-refresh-key-change-this-in-production
JWT_EXPIRATION=24h
JWT_REFRESH_EXPIRATION=168h
```

**⚠️ PENTING:** Ganti `JWT_SECRET` dan `JWT_REFRESH_SECRET` dengan key yang aman untuk production!

---

## 🚀 Menjalankan Aplikasi

### Development Mode

#### Cara 1: Direct Run
```bash
go run main.go
```

#### Cara 2: Build & Run
```bash
# Build
go build -o app.exe .

# Run
./app.exe
```

#### Cara 3: Live Reload (Recommended)
```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with Air
air
```

### Production Mode

```bash
# Build dengan optimisasi
go build -ldflags="-s -w" -o app.exe .

# Set environment
$env:APP_ENV="production"

# Run
./app.exe
```

### Menggunakan PowerShell Script

```powershell
# Start server
.\app.exe

# Stop server (Ctrl+C atau)
Stop-Process -Name "app"
```

---

## 🧪 Testing

### Health Check

```bash
curl http://localhost:3000/health
```

**Expected Response:**
```json
{
  "status": "success",
  "message": "Server is running"
}
```

### Test Login

```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

### PowerShell Testing

```powershell
# Login as Admin
$response = Invoke-RestMethod -Uri "http://localhost:3000/api/v1/auth/login" `
  -Method POST `
  -Body (@{username="admin"; password="admin123"} | ConvertTo-Json) `
  -ContentType "application/json"

# Get token
$token = $response.data.token
$headers = @{Authorization = "Bearer $token"}

# Test endpoint
Invoke-RestMethod -Uri "http://localhost:3000/api/v1/users" `
  -Headers $headers `
  -Method GET
```

---

## 📚 Dokumentasi API

Dokumentasi lengkap API tersedia di file:
👉 **[API_DOCUMENTATION.md](./API_DOCUMENTATION.md)**

### Quick Links

- [Authentication Endpoints](./API_DOCUMENTATION.md#authentication)
- [User Management](./API_DOCUMENTATION.md#user-management)
- [Achievement Operations](./API_DOCUMENTATION.md#achievement-operations)
- [Verification Workflow](./API_DOCUMENTATION.md#verification-workflow)
- [Reports & Analytics](./API_DOCUMENTATION.md#reports--analytics)

### Base URL
```
http://localhost:3000/api/v1
```

### Authentication
Semua protected endpoints memerlukan header:
```
Authorization: Bearer <your-jwt-token>
```

---

## 👥 Role & Permission

### Roles

| Role | Description |
|------|-------------|
| **Admin** | Full system access |
| **Mahasiswa** | Student user |
| **Dosen Wali** | Academic advisor |

### Permissions

| Permission | Description | Roles |
|-----------|-------------|-------|
| `user:manage` | Full user management | Admin |
| `user:create` | Create users | Admin |
| `user:read` | View users | Admin, Dosen Wali |
| `user:update` | Update users | Admin |
| `user:delete` | Delete users | Admin |
| `achievement:create` | Create achievements | Mahasiswa |
| `achievement:read` | View achievements | All |
| `achievement:update` | Update achievements | Mahasiswa (owner) |
| `achievement:delete` | Delete achievements | Mahasiswa (owner), Admin |
| `achievement:verify` | Verify achievements | Dosen Wali, Admin |
| `report:read` | View reports | Admin, Dosen Wali |

---

## 🗄 Struktur Database

### PostgreSQL Tables

#### users
```sql
id              UUID PRIMARY KEY
username        VARCHAR(50) UNIQUE
email           VARCHAR(100) UNIQUE
password_hash   VARCHAR(255)
full_name       VARCHAR(100)
is_active       BOOLEAN
role_id         UUID FOREIGN KEY → roles(id)
created_at      TIMESTAMP
updated_at      TIMESTAMP
deleted_at      TIMESTAMP (soft delete)
```

#### roles
```sql
id              UUID PRIMARY KEY
name            VARCHAR(50) UNIQUE
description     TEXT
created_at      TIMESTAMP
```

#### permissions
```sql
id              UUID PRIMARY KEY
name            VARCHAR(100) UNIQUE
description     TEXT
created_at      TIMESTAMP
```

#### role_permissions
```sql
role_id         UUID FOREIGN KEY → roles(id)
permission_id   UUID FOREIGN KEY → permissions(id)
PRIMARY KEY (role_id, permission_id)
```

#### students
```sql
id              UUID PRIMARY KEY
user_id         UUID UNIQUE FOREIGN KEY → users(id)
student_id      VARCHAR(20) UNIQUE
program_study   VARCHAR(100)
academic_year   VARCHAR(10)
advisor_id      UUID FOREIGN KEY → lecturers(id)
created_at      TIMESTAMP
```

#### lecturers
```sql
id              UUID PRIMARY KEY
user_id         UUID UNIQUE FOREIGN KEY → users(id)
lecturer_id     VARCHAR(20) UNIQUE
department      VARCHAR(100)
created_at      TIMESTAMP
```

### MongoDB Collections

#### achievements
```javascript
{
  _id: ObjectId,
  student_id: UUID,
  title: String,
  description: String,
  achieved_date: ISODate,
  status: String, // draft, pending_verification, verified, rejected
  data: Object, // Flexible structure
  verification: {
    verified_by: UUID,
    verified_at: ISODate,
    comments: String
  },
  rejection: {
    rejected_by: UUID,
    rejected_at: ISODate,
    reason: String
  },
  created_at: ISODate,
  updated_at: ISODate
}
```

---

## 🔐 Default Credentials

### Admin Account
```
Username: admin
Password: admin123
Email: admin@university.ac.id
```

### Student Account
```
Username: student001
Password: student123
Email: student001@university.ac.id
Student ID: STU001
```

### Lecturer Account
```
Username: lecturer001
Password: lecturer123
Email: lecturer001@university.ac.id
Lecturer ID: LEC001
```

**⚠️ PENTING:** Ganti password default untuk production!

---

## 📊 Status Codes

| Code | Description |
|------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Resource created |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Authentication required |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error |

---

## 🔄 Achievement Status Flow

```
[Draft] → [Pending Verification] → [Verified]
                                  ↘ [Rejected] → (Can resubmit)
```

1. **Draft** - Baru dibuat, belum submit
2. **Pending Verification** - Sudah disubmit, menunggu verifikasi
3. **Verified** - Disetujui oleh dosen wali
4. **Rejected** - Ditolak, bisa diperbaiki dan submit ulang

---

## 🛡 Security Features

- ✅ Password hashing dengan bcrypt (cost 10)
- ✅ JWT token dengan expiration
- ✅ Refresh token untuk token renewal
- ✅ Role-Based Access Control (RBAC)
- ✅ Permission checking middleware
- ✅ SQL injection protection (GORM)
- ✅ Soft delete untuk data integrity
- ✅ CORS configuration
- ✅ Rate limiting (dapat dikonfigurasi)

---

## 🚧 Troubleshooting

### Database Connection Error

**Problem:** `connection refused` atau `authentication failed`

**Solution:**
1. Pastikan PostgreSQL dan MongoDB berjalan
2. Cek credentials di `.env`
3. Test connection manual:
   ```bash
   psql -U postgres -d achievement_db
   mongosh
   ```

### Port Already in Use

**Problem:** `address already in use`

**Solution:**
```powershell
# Cari process di port 3000
Get-Process -Id (Get-NetTCPConnection -LocalPort 3000).OwningProcess

# Kill process
Stop-Process -Id <PID>

# Atau ganti port di .env
APP_PORT=3001
```

### JWT Token Invalid

**Problem:** `401 Unauthorized` atau `invalid token`

**Solution:**
1. Pastikan token belum expired
2. Cek format header: `Authorization: Bearer <token>`
3. Request token baru via `/auth/refresh`

### Permission Denied

**Problem:** `403 Forbidden`

**Solution:**
1. Cek role user yang login
2. Pastikan permission sesuai di database
3. Review permission matrix

---

## 📝 Development Notes

### Adding New Endpoint

1. **Define route** di `routes/routes.go`
2. **Create service method** di `service/`
3. **Add repository method** (jika perlu) di `repository/`
4. **Update documentation** di `API_DOCUMENTATION.md`

### Adding New Permission

1. **Add to seed** di `database/seed.go`
2. **Assign to role** via `role_permissions`
3. **Add middleware** di route dengan `RequirePermission()`

---

## 🤝 Kontributor

- **Ahmad Lazim** - Initial work - [Ahmadlazim-03](https://github.com/Ahmadlazim-03)

---

## 📄 License

This project is licensed under the MIT License.

---

## 📞 Support

Jika ada pertanyaan atau issue:
- Buka issue di [GitHub Issues](https://github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend/issues)
- Email: [your-email@university.ac.id]

---

**Built with ❤️ using Go & Fiber**
