# Panduan Deploy ke Vercel

## 🚀 Langkah-langkah Deployment

### 1. Persiapan

Pastikan Anda sudah memiliki:
- Account Vercel (https://vercel.com)
- Database PostgreSQL dan MongoDB sudah running di Railway
- Environment variables dari Railway

### 2. Install Vercel CLI (Sudah dilakukan)

```bash
npm install -g vercel
```

### 3. Login ke Vercel

```bash
vercel login
```

Pilih metode login (GitHub, GitLab, Bitbucket, atau Email)

### 4. Deploy ke Vercel

#### Deploy untuk Preview (Development)
```bash
vercel
```

#### Deploy ke Production
```bash
vercel --prod
```

### 5. Set Environment Variables di Vercel

Setelah deploy, Anda perlu set environment variables di Vercel Dashboard:

1. Buka project di Vercel Dashboard
2. Pergi ke **Settings** → **Environment Variables**
3. Tambahkan variable berikut:

```
DATABASE_URL=postgresql://username:password@host:port/database
MONGO_URI=mongodb://username:password@host:port/database
JWT_SECRET=your-super-secret-jwt-key-change-this
PORT=3000
```

**Dapatkan nilai-nilai ini dari Railway:**
- Login ke Railway (https://railway.app)
- Buka PostgreSQL service → Variables tab → Copy `DATABASE_URL`
- Buka MongoDB service → Variables tab → Copy `MONGO_URI`

### 6. Redeploy Setelah Set Environment Variables

Setelah menambahkan environment variables, redeploy:

```bash
vercel --prod
```

### 7. Test Deployment

Setelah deploy selesai, Vercel akan memberikan URL seperti:
```
https://your-project-name.vercel.app
```

Test API endpoint:
```bash
curl https://your-project-name.vercel.app/api/v1/health
```

Test Frontend:
```
https://your-project-name.vercel.app
```

## 📋 Struktur Deployment di Vercel

```
Frontend (Static)          → https://your-project.vercel.app
Backend API (Serverless)   → https://your-project.vercel.app/api/v1/*
```

## ✅ Checklist Deployment

- [ ] Login ke Vercel CLI
- [ ] Deploy dengan `vercel --prod`
- [ ] Set DATABASE_URL dari Railway
- [ ] Set MONGO_URI dari Railway
- [ ] Set JWT_SECRET
- [ ] Redeploy setelah set environment variables
- [ ] Test API endpoint `/api/v1/health`
- [ ] Test login page
- [ ] Test register user baru
- [ ] Test dashboard

## 🔧 Troubleshooting

### Database Connection Failed
- Pastikan `DATABASE_URL` dan `MONGO_URI` sudah benar
- Railway database harus allow external connections
- Check whitelist IP di Railway (set to `0.0.0.0/0` untuk allow semua)

### API Not Found (404)
- Pastikan file `api/index.go` ada
- Check vercel.json configuration
- Lihat logs di Vercel Dashboard → Deployments → View Function Logs

### Build Failed
- Check build logs di Vercel Dashboard
- Pastikan `go.mod` dan dependencies sudah benar
- Test build lokal: `cd frontend && npm run build`

## 📝 Important Notes

1. **Cold Start**: Serverless functions mungkin agak lambat di request pertama (cold start)
2. **Database**: Gunakan connection pooling untuk better performance
3. **Logs**: Check logs di Vercel Dashboard untuk debugging
4. **Custom Domain**: Bisa add custom domain di Vercel Settings
5. **Auto Deploy**: Connect ke GitHub untuk auto deploy on push

## 🔗 Useful Links

- Vercel Dashboard: https://vercel.com/dashboard
- Vercel Docs: https://vercel.com/docs
- Railway Dashboard: https://railway.app/dashboard
- Go on Vercel: https://vercel.com/docs/functions/serverless-functions/runtimes/go

## 🎯 Next Steps Setelah Deploy

1. Test semua fitur di production
2. Setup custom domain (optional)
3. Enable GitHub auto-deployment (optional)
4. Setup monitoring dan alerts
5. Configure CORS jika perlu untuk specific domains
