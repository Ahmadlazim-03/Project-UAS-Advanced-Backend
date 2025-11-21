# 🚀 Quick Deploy to Vercel

Panduan singkat deploy Achievement System ke Vercel dalam 5 menit!

## Prerequisites

- ✅ Akun GitHub
- ✅ Akun Vercel ([sign up gratis](https://vercel.com))
- ✅ Akun Railway ([sign up gratis](https://railway.app)) untuk database

---

## Step 1: Setup Database (2 menit)

### Railway - PostgreSQL

1. Login ke [railway.app](https://railway.app)
2. **New Project** → **Provision PostgreSQL**
3. Klik database → **Variables** tab
4. Copy nilai dari:
   - `PGHOST` (DB_HOST)
   - `PGUSER` (DB_USER)
   - `PGPASSWORD` (DB_PASSWORD)
   - `PGDATABASE` (DB_NAME)
   - `PGPORT` (DB_PORT)

### Railway - MongoDB

1. Di project yang sama → **New** → **Database** → **Add MongoDB**
2. Copy **MONGO_URL** dari Variables tab

---

## Step 2: Deploy ke Vercel (2 menit)

### Via Vercel Dashboard

1. **Login** ke [vercel.com/dashboard](https://vercel.com/dashboard)

2. **Import Project**
   - Click "Add New..." → "Project"
   - Import dari GitHub
   - Pilih repository ini

3. **Deploy** (klik Deploy, akan fail - itu normal!)

4. **Add Environment Variables**
   - Go to **Settings** → **Environment Variables**
   - Tambahkan variable ini (dari Railway):
   
   ```
   Name: DB_HOST
   Value: containers-us-west-xxx.railway.app
   
   Name: DB_USER  
   Value: postgres
   
   Name: DB_PASSWORD
   Value: [password dari Railway]
   
   Name: DB_NAME
   Value: railway
   
   Name: DB_PORT
   Value: 5432
   
   Name: MONGO_URI
   Value: mongodb://mongo:[password]@[host]:27017/achievement_db
   
   Name: JWT_SECRET
   Value: buatSecretKeyRandomMinimum32CharactersUntukProduction123
   
   Name: PORT
   Value: 3000
   
   Name: GIN_MODE
   Value: release
   ```

5. **Redeploy**
   - **Deployments** tab
   - Click "..." pada deployment terakhir
   - Click **"Redeploy"**
   - ✅ Done!

---

## Step 3: Test (1 menit)

1. **Buka** `https://your-project.vercel.app`
2. **Test API**: `https://your-project.vercel.app/api/v1/health`
3. **Login** dengan:
   - Username: `admin`
   - Password: `admin123`

---

## 🎉 Selesai!

Website Anda sudah live di Vercel!

### Troubleshooting

**Database connection error?**
- Pastikan semua environment variables benar
- Check Railway database masih running
- Format MONGO_URI: `mongodb://user:pass@host:port/dbname`

**Build failed?**
- Clear build cache di Settings → General
- Redeploy

**API 404?**
- Check Function Logs di Vercel
- Pastikan `vercel.json` ada di root folder

---

## 📚 More Info

- **Full Guide**: [DEPLOYMENT.md](./DEPLOYMENT.md)
- **API Docs**: `/api/v1/health`
- **Frontend Docs**: [frontend/FRONTEND_README.md](./frontend/FRONTEND_README.md)

---

**Need Help?** Create an issue di GitHub!
