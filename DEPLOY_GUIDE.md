# 🚀 Quick Deploy Guide - Vercel & Railway

## Prerequisites

- GitHub account
- Vercel account (free tier works)
- Railway account (for database)
- PostgreSQL & MongoDB databases already set up on Railway

## 📋 Deployment Steps

### Step 1: Prepare Database Connections

1. Go to your Railway dashboard
2. Copy your PostgreSQL connection string
3. Copy your MongoDB connection string
4. Keep these handy for later

### Step 2: Deploy Backend to Railway

```bash
# Login to Railway
railway login

# Initialize in project root
cd /workspaces/Project-UAS-Advanced-Backend
railway init

# Link to your project or create new
railway link

# Add environment variables
railway variables set PORT=3000
railway variables set DB_HOST=your-postgres-host
railway variables set DB_USER=your-db-user
railway variables set DB_PASSWORD=your-db-password
railway variables set DB_NAME=achievement_db
railway variables set JWT_SECRET=your-secret-key-here
railway variables set MONGO_URI=your-mongodb-connection-string

# Deploy
railway up
```

### Step 3: Get Backend URL

After deployment, get your backend URL:
```bash
railway domain
```

You'll get something like: `https://your-app.railway.app`

### Step 4: Update Frontend API URL

Edit `frontend/src/lib/api.ts`:

```typescript
const API_URL = 'https://your-app.railway.app/api/v1';
```

### Step 5: Deploy Frontend to Vercel

```bash
# Install Vercel CLI
npm i -g vercel

# Navigate to frontend
cd frontend

# Deploy
vercel

# Follow the prompts:
# - Set up and deploy? Yes
# - Which scope? Your account
# - Link to existing project? No
# - Project name? achievement-frontend
# - Directory? ./
# - Build settings? No

# After first deploy, deploy to production
vercel --prod
```

### Step 6: Configure CORS (if needed)

If you get CORS errors, update `main.go`:

```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "https://your-vercel-app.vercel.app",
    AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
    AllowHeaders: "Origin, Content-Type, Accept, Authorization",
}))
```

Then redeploy backend:
```bash
railway up
```

## 🎯 One-Command Deploy

For future deployments:

### Backend:
```bash
railway up
```

### Frontend:
```bash
cd frontend && vercel --prod
```

## 🔄 Environment Variables

### Backend (Railway)
- `PORT` - 3000
- `DB_HOST` - Your PostgreSQL host
- `DB_USER` - Your PostgreSQL user
- `DB_PASSWORD` - Your PostgreSQL password
- `DB_NAME` - achievement_db
- `JWT_SECRET` - Random secret string
- `MONGO_URI` - MongoDB connection string

### Frontend (Vercel)
No environment variables needed if you hardcode the API URL.

Alternatively, use Vercel environment variable:
- `PUBLIC_API_URL` - Your Railway backend URL

## ✅ Testing

1. **Backend**: Visit `https://your-app.railway.app/api/v1`
   - Should return: "Student Achievement System API"

2. **Frontend**: Visit your Vercel URL
   - Should show the login page

3. **Full Test**: Register and login
   - Create achievement
   - Test all features

## 🐛 Troubleshooting

### Database Connection Failed
- Check Railway database is running
- Verify connection strings
- Check firewall rules

### CORS Errors
- Update CORS config in `main.go`
- Redeploy backend

### Frontend 404 Errors
- Check API_URL in `frontend/src/lib/api.ts`
- Ensure backend is deployed and running

### Build Failures
- Check logs: `railway logs` or Vercel dashboard
- Ensure all dependencies are listed

## 📱 Access Your App

- **Frontend**: `https://your-app.vercel.app`
- **Backend API**: `https://your-backend.railway.app`
- **Swagger Docs**: `https://your-backend.railway.app/swagger/index.html`

## 🔐 Default Credentials

After seeding, you can login with:
- **Admin**: username: admin, password: admin123
- **Dosen**: username: dosen1, password: dosen123  
- **Mahasiswa**: username: student1, password: student123

## 📝 Notes

- Railway free tier gives 500 hours/month
- Vercel free tier is generous for hobby projects
- Database restarts if idle for 30 days on Railway
- Always use environment variables for secrets
