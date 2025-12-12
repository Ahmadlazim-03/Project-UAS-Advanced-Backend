# 🎓 Student Achievement System - Frontend

Modern web application built with React + Vite for managing student achievements.

## 🚀 Features

- ✅ **Authentication** - Login & Register with JWT
- ✅ **Role-Based Dashboards** - Admin, Student (Mahasiswa), Lecturer (Dosen Wali)
- ✅ **Achievement Management** - CRUD operations for achievements
- ✅ **Verification Workflow** - Submit, verify, and reject achievements
- ✅ **Responsive Design** - Works on desktop and mobile
- ✅ **Modern UI** - Tailwind CSS with custom components

## 🛠️ Tech Stack

- **React 18** - UI library
- **Vite** - Build tool & dev server
- **React Router** - Client-side routing
- **Axios** - HTTP client
- **Tailwind CSS** - Utility-first CSS framework
- **Lucide React** - Icon library

## 📦 Installation

### Prerequisites

- Node.js 18+ 
- npm or yarn
- Backend API running on `http://localhost:3000`

### Setup

1. **Install dependencies**
```bash
npm install
```

2. **Configure environment**
```bash
# .env file is already created with default values
# Edit if your backend runs on a different port
VITE_API_BASE_URL=http://localhost:3000/api/v1
```

3. **Start development server**
```bash
npm run dev
```

The app will run on `http://localhost:5173`

## 📁 Project Structure

```
App/
├── src/
│   ├── components/        # Reusable components
│   │   ├── DashboardLayout.jsx
│   │   └── Loading.jsx
│   ├── contexts/          # React contexts
│   │   └── AuthContext.jsx
│   ├── pages/             # Page components
│   │   ├── Login.jsx
│   │   ├── Register.jsx
│   │   ├── Admin/
│   │   │   ├── Dashboard.jsx
│   │   │   ├── Users.jsx
│   │   │   └── Reports.jsx
│   │   ├── Student/
│   │   │   ├── Dashboard.jsx
│   │   │   └── Achievements.jsx
│   │   └── Lecturer/
│   │       ├── Dashboard.jsx
│   │       └── Achievements.jsx
│   ├── services/          # API services
│   │   ├── api.js
│   │   └── index.js
│   ├── utils/             # Helper functions
│   │   ├── auth.js
│   │   └── helpers.js
│   ├── App.jsx            # Main app component
│   ├── main.jsx           # Entry point
│   └── index.css          # Global styles
├── package.json
├── vite.config.js
├── tailwind.config.js
└── README.md
```

## 🔐 Authentication

### Default Accounts

After backend seeding, you can login with:

| Role | Username | Password |
|------|----------|----------|
| Admin | `admin` | `admin123` |
| Student | `student001` | `student123` |
| Lecturer | `lecturer001` | `lecturer123` |

## 🎨 Available Scripts

```bash
# Development
npm run dev          # Start dev server (port 5173)

# Build
npm run build        # Build for production

# Preview
npm run preview      # Preview production build

# Lint
npm run lint         # Run ESLint
```

## 📱 Features by Role

### Admin
- View system statistics
- Manage users (CRUD)
- Access reports and analytics
- Monitor all achievements

### Student (Mahasiswa)
- Create and manage achievements
- Submit for verification
- Track achievement status
- View personal statistics

### Lecturer (Dosen Wali)
- View advisee achievements
- Verify or reject submissions
- Add verification comments
- Track advisee progress

## 🔌 API Integration

The app connects to the backend API:
- Base URL: `http://localhost:3000/api/v1`
- Authentication: JWT Bearer tokens
- Automatic token refresh
- Error handling with interceptors

## 🎯 Key Components

### DashboardLayout
Responsive sidebar layout with:
- Role-based navigation
- User profile display
- Logout functionality
- Mobile-friendly menu

### AuthContext
Global authentication state:
- Login/logout functions
- User data management
- Protected route handling
- Token storage

### Protected Routes
Automatic role-based access control:
- Admin-only routes
- Student-only routes
- Lecturer-only routes
- Redirect unauthorized users

## 🚧 Development

### Adding New Pages

1. Create component in `src/pages/[Role]/`
2. Add route in `src/App.jsx`
3. Update navigation in `DashboardLayout.jsx`

### Adding New API Calls

1. Add service function in `src/services/index.js`
2. Use in components with try-catch
3. Handle loading and error states

## 📝 Environment Variables

```env
VITE_API_BASE_URL=http://localhost:3000/api/v1
```

## 🔒 Security

- JWT token stored in localStorage
- Automatic token expiry handling
- Protected routes with role checking
- CORS configured in backend

## 🐛 Troubleshooting

### Backend Connection Error
```
Error: Network Error
```
**Solution:** Ensure backend is running on port 3000

### Login Failed
```
401 Unauthorized
```
**Solution:** Check credentials or clear localStorage

### Build Errors
```bash
# Clear node_modules and reinstall
rm -rf node_modules
npm install
```

## 📄 License

MIT License

## 👨‍💻 Development Team

Built for UAS Advanced Backend Development Project

---

**Happy Coding! 🚀**
