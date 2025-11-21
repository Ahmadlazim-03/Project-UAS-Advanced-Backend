# 🚀 Quick Start Commands

## Development

### Backend
```bash
go run main.go
```

### Frontend
```bash
cd frontend && npm run dev
```

## Build

### Frontend Only
```bash
cd frontend && npm run build
```

### Full Project
```bash
./build.sh
```

## Deploy

### Backend to Railway
```bash
railway up
```

### Frontend to Vercel
```bash
cd frontend && vercel --prod
```

## URLs

- **Local Backend**: http://localhost:3000
- **Local Frontend**: http://localhost:5173
- **Swagger Docs**: http://localhost:3000/swagger/index.html

## Default Login

**Admin:**
- Username: `admin`
- Password: `admin123`

**Dosen Wali:**
- Username: `dosen1`
- Password: `dosen123`

**Mahasiswa:**
- Username: `student1`
- Password: `student123`

## Environment Variables

Create `.env` in root:
```env
PORT=3000
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=achievement_db
JWT_SECRET=your-secret-key-here
MONGO_URI=mongodb://localhost:27017/achievements
```

## Helpful Commands

```bash
# Check backend is running
curl http://localhost:3000/api/v1

# Check database connection
psql -h localhost -U postgres -d achievement_db

# View Railway logs
railway logs

# View Vercel deployment logs
vercel logs

# Update frontend build
cd frontend && npm run build && cd ..

# Test backend
go test ./...

# Format Go code
go fmt ./...

# Clean builds
rm -rf frontend/build frontend/.svelte-kit
```

## Troubleshooting

**Backend not starting?**
- Check .env file exists
- Verify database connections
- Check port 3000 is available

**Frontend build errors?**
- Run `npm install` in frontend/
- Clear .svelte-kit and build folders
- Check Node.js version (need 18+)

**CORS errors?**
- Update CORS config in main.go
- Check API URL in frontend/src/lib/api.ts

**Database errors?**
- Verify database is running
- Check connection strings
- Run migrations if needed

## Documentation

- `README.md` - Main documentation
- `DEPLOY_GUIDE.md` - Deployment instructions
- `DEPLOYMENT.md` - Advanced deployment options
- `PROJECT_STATUS.md` - Current project status
- `frontend/FRONTEND_README.md` - Frontend specific docs

---

**Ready to deploy? See `DEPLOY_GUIDE.md`**
