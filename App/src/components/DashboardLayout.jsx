import { Link, useNavigate, useLocation } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { 
  LayoutDashboard, 
  Users, 
  Award, 
  FileText, 
  LogOut, 
  Menu, 
  X,
  GraduationCap,
  UserCheck,
  Trash2
} from 'lucide-react'
import { useState } from 'react'

export default function DashboardLayout({ children, title }) {
  const { user, logout } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const [sidebarOpen, setSidebarOpen] = useState(false)

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  const getMenuItems = () => {
    if (user?.role_name === 'Admin') {
      return [
        { icon: LayoutDashboard, label: 'Dashboard', path: '/admin/dashboard' },
        { icon: Award, label: 'Achievements', path: '/admin/achievements' },
        { icon: Users, label: 'Users', path: '/admin/users' },
        { icon: UserCheck, label: 'Advisors', path: '/admin/advisors' },
        { icon: Trash2, label: 'Deleted Users', path: '/admin/deleted-users' },
        { icon: FileText, label: 'Reports', path: '/admin/reports' },
      ]
    } else if (user?.role_name === 'Mahasiswa') {
      return [
        { icon: LayoutDashboard, label: 'Dashboard', path: '/student/dashboard' },
        { icon: Award, label: 'My Achievements', path: '/student/achievements' },
      ]
    } else if (user?.role_name === 'Dosen Wali') {
      return [
        { icon: LayoutDashboard, label: 'Dashboard', path: '/lecturer/dashboard' },
        { icon: Users, label: 'My Advisees', path: '/lecturer/advisees' },
        { icon: Award, label: 'Student Achievements', path: '/lecturer/achievements' },
      ]
    }
    return []
  }

  const menuItems = getMenuItems()

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Sidebar */}
      <aside
        className={`fixed inset-y-0 left-0 z-50 w-64 bg-white shadow-lg transform ${
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        } lg:translate-x-0 transition-transform duration-300 ease-in-out`}
      >
        <div className="flex items-center justify-between h-16 px-6 border-b">
          <div className="flex items-center space-x-2">
            <GraduationCap className="w-8 h-8 text-primary-600" />
            <span className="text-xl font-bold text-gray-800">Achievement</span>
          </div>
          <button
            onClick={() => setSidebarOpen(false)}
            className="lg:hidden text-gray-500 hover:text-gray-700"
          >
            <X className="w-6 h-6" />
          </button>
        </div>

        <nav className="px-4 py-6 space-y-2">
          {menuItems.map((item) => {
            const Icon = item.icon
            const isActive = location.pathname === item.path
            return (
              <Link
                key={item.path}
                to={item.path}
                className={`flex items-center space-x-3 px-4 py-3 rounded-lg transition-colors ${
                  isActive
                    ? 'bg-primary-50 text-primary-600'
                    : 'text-gray-700 hover:bg-gray-100'
                }`}
              >
                <Icon className="w-5 h-5" />
                <span className="font-medium">{item.label}</span>
              </Link>
            )
          })}
        </nav>

        <div className="absolute bottom-0 left-0 right-0 p-4 border-t">
          <div className="px-4 py-3 mb-2 bg-gray-50 rounded-lg">
            <div className="flex items-center space-x-3">
              <div className="flex items-center justify-center w-10 h-10 rounded-full bg-primary-100">
                <UserCheck className="w-5 h-5 text-primary-600" />
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-gray-900 truncate">
                  {user?.full_name || user?.username}
                </p>
                <p className="text-xs text-gray-500">{user?.role_name}</p>
              </div>
            </div>
          </div>
          <button
            onClick={handleLogout}
            className="flex items-center space-x-3 w-full px-4 py-3 text-gray-700 hover:bg-red-50 hover:text-red-600 rounded-lg transition-colors"
          >
            <LogOut className="w-5 h-5" />
            <span className="font-medium">Logout</span>
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <div className="lg:pl-64">
        {/* Header */}
        <header className="bg-white shadow-sm sticky top-0 z-40">
          <div className="flex items-center justify-between h-16 px-6">
            <div className="flex items-center space-x-4">
              <button
                onClick={() => setSidebarOpen(true)}
                className="lg:hidden text-gray-500 hover:text-gray-700"
              >
                <Menu className="w-6 h-6" />
              </button>
              <h1 className="text-2xl font-bold text-gray-800">{title}</h1>
            </div>
            <div className="flex items-center space-x-4">
              <span className="text-sm text-gray-600">{user?.email}</span>
            </div>
          </div>
        </header>

        {/* Page Content */}
        <main className="p-6">{children}</main>
      </div>

      {/* Overlay */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 bg-black bg-opacity-50 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}
    </div>
  )
}
