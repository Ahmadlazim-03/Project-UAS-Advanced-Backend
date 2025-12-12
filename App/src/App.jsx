import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { useAuth } from './contexts/AuthContext'
import Login from './pages/Login'
import Register from './pages/Register'
import AdminDashboard from './pages/Admin/Dashboard'
import AdminAchievements from './pages/Admin/Achievements'
import StudentDashboard from './pages/Student/Dashboard'
import LecturerDashboard from './pages/Lecturer/Dashboard'
import StudentAchievements from './pages/Student/Achievements'
import LecturerAchievements from './pages/Lecturer/Achievements'
import Reports from './pages/Admin/Reports'
import Users from './pages/Admin/Users'
import Advisors from './pages/Admin/Advisors'
import DeletedUsers from './pages/Admin/DeletedUsers'
import Loading from './components/Loading'
import { useMemo } from 'react'

function ProtectedRoute({ children, allowedRoles }) {
  const { user, loading } = useAuth()

  if (loading) {
    return <Loading />
  }

  if (!user) {
    return <Navigate to="/login" replace />
  }

  if (allowedRoles && !allowedRoles.includes(user.role_name)) {
    return <Navigate to="/unauthorized" replace />
  }

  return children
}

function DashboardRedirect() {
  const { user, loading } = useAuth()

  const redirectPath = useMemo(() => {
    if (!user) return '/login'
    
    switch (user.role_name) {
      case 'Admin':
        return '/admin/dashboard'
      case 'Mahasiswa':
        return '/student/dashboard'
      case 'Dosen Wali':
        return '/lecturer/dashboard'
      default:
        return '/login'
    }
  }, [user])

  if (loading) {
    return <Loading />
  }

  return <Navigate to={redirectPath} replace />
}

function App() {
  return (
    <Router>
      <Routes>
        {/* Public Routes */}
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

        {/* Protected Routes - Admin */}
        <Route
          path="/admin/dashboard"
          element={
            <ProtectedRoute allowedRoles={['Admin']}>
              <AdminDashboard />
            </ProtectedRoute>
          }
        />
        <Route
          path="/admin/users"
          element={
            <ProtectedRoute allowedRoles={['Admin']}>
              <Users />
            </ProtectedRoute>
          }
        />
        <Route
          path="/admin/advisors"
          element={
            <ProtectedRoute allowedRoles={['Admin']}>
              <Advisors />
            </ProtectedRoute>
          }
        />
        <Route
          path="/admin/deleted-users"
          element={
            <ProtectedRoute allowedRoles={['Admin']}>
              <DeletedUsers />
            </ProtectedRoute>
          }
        />
        <Route
          path="/admin/reports"
          element={
            <ProtectedRoute allowedRoles={['Admin']}>
              <Reports />
            </ProtectedRoute>
          }
        />
        <Route
          path="/admin/achievements"
          element={
            <ProtectedRoute allowedRoles={['Admin']}>
              <AdminAchievements />
            </ProtectedRoute>
          }
        />

        {/* Protected Routes - Student */}
        <Route
          path="/student/dashboard"
          element={
            <ProtectedRoute allowedRoles={['Mahasiswa']}>
              <StudentDashboard />
            </ProtectedRoute>
          }
        />
        <Route
          path="/student/achievements"
          element={
            <ProtectedRoute allowedRoles={['Mahasiswa']}>
              <StudentAchievements />
            </ProtectedRoute>
          }
        />

        {/* Protected Routes - Lecturer */}
        <Route
          path="/lecturer/dashboard"
          element={
            <ProtectedRoute allowedRoles={['Dosen Wali']}>
              <LecturerDashboard />
            </ProtectedRoute>
          }
        />
        <Route
          path="/lecturer/achievements"
          element={
            <ProtectedRoute allowedRoles={['Dosen Wali']}>
              <LecturerAchievements />
            </ProtectedRoute>
          }
        />

        {/* Redirect to role-based dashboard */}
        <Route path="/dashboard" element={<DashboardRedirect />} />

        {/* Default Route */}
        <Route path="/" element={<DashboardRedirect />} />
        <Route path="/unauthorized" element={<div className="flex items-center justify-center min-h-screen"><h1 className="text-2xl font-bold text-red-600">Unauthorized Access</h1></div>} />
        <Route path="*" element={<div className="flex items-center justify-center min-h-screen"><h1 className="text-2xl font-bold">404 - Page Not Found</h1></div>} />
      </Routes>
    </Router>
  )
}

export default App
