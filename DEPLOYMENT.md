# 🚀 Deployment Guide - Vercel

Panduan lengkap untuk deploy fullstack aplikasi Achievement System ke Vercel dengan database Railway.

## 📋 Prerequisites

1. **Akun Vercel** - [Sign up di vercel.com](https://vercel.com)
2. **Akun Railway** - [Sign up di railway.app](https://railway.app) untuk database
3. **Akun GitHub** - Repository sudah ter-push ke GitHub
4. **Vercel CLI** (optional) - `npm i -g vercel`

---

## 🗄️ Step 1: Setup Database di Railway

### PostgreSQL Database

1. Login ke [Railway.app](https://railway.app)
2. Click **"New Project"**
3. Select **"Provision PostgreSQL"**
4. Tunggu database siap
5. Click database → **"Connect"** → Copy connection details:
   - `DB_HOST`
   - `DB_USER` 
   - `DB_PASSWORD`
   - `DB_NAME`
   - `DB_PORT`

### MongoDB Database

1. Di Railway project yang sama, click **"New"**
2. Select **"Database"** → **"Add MongoDB"**
3. Tunggu database siap
4. Copy `MONGO_URI` dari connection string

---

## 🔧 Step 2: Konfigurasi Project

### Update API URL di Frontend

Edit `frontend/src/lib/api.ts`:

```typescript
const API_URL = '/api/v1'; // Sudah benar untuk Vercel
```

File ini sudah dikonfigurasi dengan benar!

---

## ☁️ Step 3: Deploy ke Vercel

### Method 1: Via Vercel Dashboard (Recommended)

1. **Login ke Vercel**
   - Buka [vercel.com/dashboard](https://vercel.com/dashboard)
   - Login dengan GitHub

2. **Import Project**
   - Click **"Add New Project"**
   - Click **"Import Git Repository"**
   - Pilih repository `Project-UAS-Advanced-Backend`
   - Click **"Import"**

3. **Configure Project**
   - **Framework Preset**: Other
   - **Root Directory**: `./` (leave as is)
   - **Build Command**: (leave empty, Vercel will auto-detect)
   - **Output Directory**: (leave empty)
   - Click **"Deploy"** (will fail first time, it's okay)

4. **Add Environment Variables**
   - Go to **Settings** → **Environment Variables**
   - Add semua variable dari `.env.example`:
   
   ```
   DB_HOST=containers-us-west-123.railway.app
   DB_USER=postgres
   DB_PASSWORD=your-password-here
   DB_NAME=railway
   DB_PORT=5432
   MONGO_URI=mongodb://mongo:password@containers-us-west-456.railway.app:27017/achievement_db
   JWT_SECRET=your-super-secret-jwt-key-minimum-32-characters-long
   PORT=3000
   GIN_MODE=release
   ```

5. **Redeploy**
   - Go to **Deployments**
   - Click **"..."** pada deployment terakhir
   - Click **"Redeploy"**
   - Wait for deployment to complete

### Method 2: Via Vercel CLI

```bash
# Install Vercel CLI
npm i -g vercel

# Login
vercel login

# Deploy
vercel

# Follow prompts:
# - Set up and deploy? Yes
# - Which scope? Select your account
# - Link to existing project? No
# - Project name? Project-UAS-Advanced-Backend
# - Directory? ./
# - Override settings? No

# Add environment variables
vercel env add DB_HOST
vercel env add DB_USER
vercel env add DB_PASSWORD
vercel env add DB_NAME
vercel env add DB_PORT
vercel env add MONGO_URI
vercel env add JWT_SECRET

# Deploy to production
vercel --prod
```

---

## ✅ Step 4: Verify Deployment

1. **Test API**
   ```bash
   curl https://your-app.vercel.app/api/v1/health
   ```
   
   Expected response:
   ```json
   {
     "status": "success",
     "message": "API is running"
   }
   ```

2. **Test Frontend**
   - Buka `https://your-app.vercel.app`
   - Anda harus melihat halaman login
   - Try login dengan kredensial yang sudah di-seed

3. **Check Logs**
   - Vercel Dashboard → Your Project → **Deployments**
   - Click deployment → **View Function Logs**
   - Check untuk errors

---

## 🔍 Troubleshooting

### Database Connection Failed

**Error**: `connection refused` atau `timeout`

**Solution**:
1. Pastikan Railway database sudah running
2. Check environment variables di Vercel benar
3. Railway PostgreSQL/MongoDB allow connections dari anywhere
4. Format MONGO_URI harus benar: `mongodb://user:pass@host:port/dbname`

### Build Failed

**Error**: Build fails saat deployment

**Solution**:
1. Check build logs di Vercel
2. Pastikan `frontend/package.json` ada `vercel-build` script
3. Pastikan adapter static terinstall: `@sveltejs/adapter-static`
4. Clear build cache: Settings → General → **Clear Build Cache**

### API 404 Errors

**Error**: `/api/v1/*` returns 404

**Solution**:
1. Check `vercel.json` routing configuration
2. Pastikan `api/index.go` ter-deploy dengan benar
3. Check Function Logs untuk errors
4. Verify rewrites di `vercel.json`

### CORS Errors

**Error**: CORS policy blocking requests

**Solution**:
API sudah dikonfigurasi dengan CORS yang permissive. Jika masih ada masalah:

Edit `api/index.go`:
```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "https://your-app.vercel.app",
    AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
    AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    AllowCredentials: true,
}))
```

### JWT Token Issues

**Error**: Authentication failing

**Solution**:
1. Pastikan `JWT_SECRET` sama di semua environments
2. Token harus minimum 32 characters
3. Check localStorage di browser (F12 → Application → Local Storage)

---

## 🔄 Update & Redeploy

### Automatic Deployments (Recommended)

Vercel akan auto-deploy setiap kali Anda push ke GitHub:

```bash
git add .
git commit -m "Update feature"
git push origin main
```

Vercel akan otomatis rebuild dan deploy!

### Manual Redeploy

Via Dashboard:
1. Vercel Dashboard → Deployments
2. Click **"..."** → **"Redeploy"**

Via CLI:
```bash
vercel --prod
```

---

## 📊 Monitoring

### View Logs

1. **Function Logs** (API):
   - Vercel Dashboard → Deployments → Click deployment
   - **Functions** tab → Click function → **Logs**

2. **Build Logs**:
   - Vercel Dashboard → Deployments → Click deployment
   - **Building** tab

### Performance

1. **Analytics**:
   - Vercel Dashboard → Analytics
   - View traffic, performance metrics

2. **Limits** (Free tier):
   - 100 GB Bandwidth
   - 6,000 Build minutes
   - Serverless Function: 10s max execution time

---

## 🎯 Production Checklist

- [ ] Database Railway sudah setup (PostgreSQL + MongoDB)
- [ ] Environment variables sudah ditambahkan di Vercel
- [ ] JWT_SECRET adalah random string (min 32 chars)
- [ ] Database credentials aman dan tidak di-commit
- [ ] Test API endpoint: `/api/v1/health`
- [ ] Test login functionality
- [ ] Test create/read/update/delete operations
- [ ] Check logs untuk errors
- [ ] Setup custom domain (optional)
- [ ] Enable HTTPS (otomatis di Vercel)

---

## 🌐 Custom Domain (Optional)

1. Vercel Dashboard → Your Project → **Settings** → **Domains**
2. Add your domain
3. Update DNS records di domain provider
4. Wait for DNS propagation (~24 hours)

---

## 💡 Tips

1. **Use Environment Variables**: Never hardcode secrets
2. **Check Logs Regularly**: Monitor for errors
3. **Database Backups**: Railway auto-backups, but verify
4. **Rate Limiting**: Implement di production
5. **Monitoring**: Setup alerts untuk downtime

---

## 📞 Support

- Vercel Docs: https://vercel.com/docs
- Railway Docs: https://docs.railway.app
- SvelteKit Docs: https://kit.svelte.dev

---

## 🎉 Done!

Your Achievement System is now live on Vercel! 🚀

URL: `https://your-project-name.vercel.app`
