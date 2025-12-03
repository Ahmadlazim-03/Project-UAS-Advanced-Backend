# Quick Start Guide

Panduan cepat untuk memulai Student Achievement System.

## 🚀 Quick Start (5 Menit)

### 1. Prerequisites Check

```powershell
# Check Go
go version  # Should be 1.21+

# Check PostgreSQL
psql --version

# Check MongoDB
mongod --version
```

### 2. Clone & Setup

```powershell
# Clone repository
git clone https://github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend.git
cd Project-UAS-Advanced-Backend

# Install dependencies
go mod download
```

### 3. Database Setup

```powershell
# Create PostgreSQL database
psql -U postgres -c "CREATE DATABASE achievement_db;"

# MongoDB akan otomatis dibuat saat aplikasi berjalan
```

### 4. Environment Setup

File `.env` sudah ada dengan konfigurasi default:

```env
APP_PORT=3000
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=admin
POSTGRES_DB=achievement_db
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=achievement_db
```

**Edit sesuai dengan konfigurasi database Anda!**

### 5. Run Application

```powershell
# Build
go build -o app.exe .

# Run
.\app.exe
```

**Server running at:** `http://localhost:3000`

---

## ✅ Verify Installation

### 1. Health Check

```powershell
Invoke-RestMethod http://localhost:3000/health
```

**Expected:**
```json
{
  "status": "success",
  "message": "Server is running"
}
```

### 2. Login Test

```powershell
$response = Invoke-RestMethod -Uri "http://localhost:3000/api/v1/auth/login" `
  -Method POST `
  -Body (@{username="admin"; password="admin123"} | ConvertTo-Json) `
  -ContentType "application/json"

$response.data.token
```

**Expected:** JWT token string

---

## 📖 Next Steps

### 1. Explore API Documentation

Buka file lengkap: **[API_DOCUMENTATION.md](./API_DOCUMENTATION.md)**

### 2. Test Endpoints

```powershell
# Login
$token = (Invoke-RestMethod -Uri "http://localhost:3000/api/v1/auth/login" `
  -Method POST `
  -Body (@{username="admin"; password="admin123"} | ConvertTo-Json) `
  -ContentType "application/json").data.token

# Set headers
$headers = @{Authorization = "Bearer $token"}

# List users
Invoke-RestMethod -Uri "http://localhost:3000/api/v1/users" -Headers $headers

# Create achievement as student
$studentToken = (Invoke-RestMethod -Uri "http://localhost:3000/api/v1/auth/login" `
  -Method POST `
  -Body (@{username="student001"; password="student123"} | ConvertTo-Json) `
  -ContentType "application/json").data.token

$studentHeaders = @{Authorization = "Bearer $studentToken"}

Invoke-RestMethod -Uri "http://localhost:3000/api/v1/achievements" `
  -Headers $studentHeaders `
  -Method POST `
  -Body (@{
    title = "My First Achievement"
    description = "Testing achievement creation"
    achieved_date = "2024-12-01"
    data = @{
      competition_name = "Test Competition"
      rank = 1
    }
  } | ConvertTo-Json) `
  -ContentType "application/json"
```

---

## 🔑 Default Accounts

### Admin
```
Username: admin
Password: admin123
```

### Student
```
Username: student001
Password: student123
```

### Lecturer
```
Username: lecturer001
Password: lecturer123
```

---

## 📝 Common Commands

### Start Server
```powershell
.\app.exe
```

### Stop Server
```powershell
# Press Ctrl+C or
Stop-Process -Name "app"
```

### Rebuild
```powershell
go build -o app.exe .
```

### Reset Database
```powershell
# Stop server first
Stop-Process -Name "app"

# Drop and recreate
$env:PGPASSWORD='admin'
psql -U postgres -d achievement_db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

# Restart server (akan auto-migrate dan seed)
.\app.exe
```

---

## 🛠 Troubleshooting

### Port 3000 Already in Use

```powershell
# Find process
Get-Process -Id (Get-NetTCPConnection -LocalPort 3000).OwningProcess

# Kill process
Stop-Process -Id <PID>

# Or change port in .env
# APP_PORT=3001
```

### Database Connection Failed

```powershell
# Test PostgreSQL
psql -U postgres -d achievement_db

# Test MongoDB
mongosh

# If failed, check:
# 1. Services are running
# 2. Credentials in .env are correct
# 3. Firewall settings
```

### "Cannot find app.exe"

```powershell
# Build first
go build -o app.exe .

# Then run
.\app.exe
```

---

## 📚 Documentation

- **README.md** - Full documentation
- **API_DOCUMENTATION.md** - Complete API reference
- **QUICK_START.md** - This file

---

## 🆘 Need Help?

- Check **[README.md](./README.md)** for detailed documentation
- Check **[API_DOCUMENTATION.md](./API_DOCUMENTATION.md)** for API reference
- Open issue on GitHub

---

**Happy Coding! 🚀**
