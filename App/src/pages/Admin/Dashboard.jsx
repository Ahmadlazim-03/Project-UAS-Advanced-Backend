import { useState, useEffect } from 'react'
import DashboardLayout from '../../components/DashboardLayout'
import { reportService } from '../../services'
import { Users, GraduationCap, Award, CheckCircle } from 'lucide-react'

export default function AdminDashboard() {
  const [stats, setStats] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchStatistics()
  }, [])

  const fetchStatistics = async () => {
    try {
      const response = await reportService.getStatistics()
      if (response.status === 'success') {
        setStats(response.data)
      }
    } catch (error) {
      console.error('Error fetching statistics:', error)
    } finally {
      setLoading(false)
    }
  }

  // Calculate total achievements from achievements object
  const getTotalAchievements = () => {
    if (!stats?.achievements) return 0
    return Object.values(stats.achievements).reduce((sum, count) => sum + count, 0)
  }

  const getVerifiedAchievements = () => {
    return stats?.achievements?.verified || 0
  }

  const statCards = [
    {
      title: 'Total Students',
      value: stats?.students || 0,
      icon: GraduationCap,
      color: 'bg-blue-500',
    },
    {
      title: 'Total Lecturers',
      value: stats?.lecturers || 0,
      icon: Users,
      color: 'bg-green-500',
    },
    {
      title: 'Total Achievements',
      value: getTotalAchievements(),
      icon: Award,
      color: 'bg-purple-500',
    },
    {
      title: 'Verified',
      value: getVerifiedAchievements(),
      icon: CheckCircle,
      color: 'bg-emerald-500',
    },
  ]

  if (loading) {
    return (
      <DashboardLayout title="Admin Dashboard">
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
        </div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout title="Admin Dashboard">
      <div className="space-y-6">
        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {statCards.map((stat, index) => {
            const Icon = stat.icon
            return (
              <div key={index} className="card">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-gray-600 mb-1">{stat.title}</p>
                    <p className="text-3xl font-bold text-gray-900">{stat.value}</p>
                  </div>
                  <div className={`${stat.color} p-3 rounded-lg`}>
                    <Icon className="w-8 h-8 text-white" />
                  </div>
                </div>
              </div>
            )
          })}
        </div>

        {/* Achievement Status Breakdown */}
        {stats && (
          <div className="card">
            <h2 className="text-xl font-bold text-gray-900 mb-4">Achievement Status</h2>
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div className="p-4 bg-blue-50 rounded-lg">
                <p className="text-sm text-gray-600">Draft</p>
                <p className="text-2xl font-bold text-blue-600">{stats.draft_achievements || 0}</p>
              </div>
              <div className="p-4 bg-yellow-50 rounded-lg">
                <p className="text-sm text-gray-600">Pending</p>
                <p className="text-2xl font-bold text-yellow-600">{stats.pending_achievements || 0}</p>
              </div>
              <div className="p-4 bg-green-50 rounded-lg">
                <p className="text-sm text-gray-600">Verified</p>
                <p className="text-2xl font-bold text-green-600">{stats.verified_achievements || 0}</p>
              </div>
              <div className="p-4 bg-red-50 rounded-lg">
                <p className="text-sm text-gray-600">Rejected</p>
                <p className="text-2xl font-bold text-red-600">{stats.rejected_achievements || 0}</p>
              </div>
            </div>
          </div>
        )}

        {/* Quick Actions */}
        <div className="card">
          <h2 className="text-xl font-bold text-gray-900 mb-4">Quick Actions</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <a href="/admin/users" className="p-4 border-2 border-gray-200 rounded-lg hover:border-primary-500 hover:bg-primary-50 transition-colors">
              <Users className="w-8 h-8 text-primary-600 mb-2" />
              <h3 className="font-semibold text-gray-900">Manage Users</h3>
              <p className="text-sm text-gray-600">Create, edit, or delete users</p>
            </a>
            <a href="/admin/reports" className="p-4 border-2 border-gray-200 rounded-lg hover:border-primary-500 hover:bg-primary-50 transition-colors">
              <Award className="w-8 h-8 text-primary-600 mb-2" />
              <h3 className="font-semibold text-gray-900">View Reports</h3>
              <p className="text-sm text-gray-600">Access detailed reports</p>
            </a>
            <div className="p-4 border-2 border-gray-200 rounded-lg hover:border-primary-500 hover:bg-primary-50 transition-colors cursor-pointer">
              <GraduationCap className="w-8 h-8 text-primary-600 mb-2" />
              <h3 className="font-semibold text-gray-900">Student Overview</h3>
              <p className="text-sm text-gray-600">View all students</p>
            </div>
          </div>
        </div>
      </div>
    </DashboardLayout>
  )
}
