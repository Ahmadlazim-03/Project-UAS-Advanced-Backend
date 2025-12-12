import { useState, useEffect } from 'react'
import DashboardLayout from '../../components/DashboardLayout'
import { lecturerService } from '../../services'
import { Users, Award, Clock, CheckCircle, XCircle } from 'lucide-react'
import { useAuth } from '../../contexts/AuthContext'

export default function LecturerDashboard() {
  const { user } = useAuth()
  const [stats, setStats] = useState({
    advisees: 0,
    pending: 0,
    verified: 0,
    rejected: 0,
    total: 0,
  })
  const [advisees, setAdvisees] = useState([])
  const [recentAchievements, setRecentAchievements] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    try {
      // Fetch advisee achievements using the correct endpoint
      const achievementsResponse = await lecturerService.getAdviseeAchievements()
      console.log('Dashboard response:', achievementsResponse)

      if (achievementsResponse.status === 'success') {
        // Backend returns paginated response: response.pagination.data.achievements
        const achievements = achievementsResponse.pagination?.data?.achievements || []
        console.log('Dashboard achievements:', achievements)
        setRecentAchievements(achievements.slice(0, 5))

        // Calculate stats from achievements (backend uses 'submitted' for pending)
        const pending = achievements.filter(a => a.status === 'submitted').length
        const verified = achievements.filter(a => a.status === 'verified').length
        const rejected = achievements.filter(a => a.status === 'rejected').length

        // Count unique students from achievements
        const uniqueStudents = new Set(achievements.map(a => a.student_id))

        setStats({
          advisees: uniqueStudents.size,
          pending,
          verified,
          rejected,
          total: achievements.length,
        })
      }
    } catch (error) {
      console.error('Error fetching data:', error)
    } finally {
      setLoading(false)
    }
  }

  const statCards = [
    {
      title: 'My Advisees',
      value: stats.advisees,
      icon: Users,
      color: 'bg-blue-500',
    },
    {
      title: 'Pending Review',
      value: stats.pending,
      icon: Clock,
      color: 'bg-yellow-500',
    },
    {
      title: 'Verified',
      value: stats.verified,
      icon: CheckCircle,
      color: 'bg-green-500',
    },
    {
      title: 'Rejected',
      value: stats.rejected,
      icon: XCircle,
      color: 'bg-red-500',
    },
  ]

  const getStatusBadge = (status) => {
    const statusMap = {
      draft: 'badge-info',
      pending_verification: 'badge-warning',
      verified: 'badge-success',
      rejected: 'badge-danger',
    }
    return statusMap[status] || 'badge-info'
  }

  if (loading) {
    return (
      <DashboardLayout title="Lecturer Dashboard">
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
        </div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout title="Lecturer Dashboard">
      <div className="space-y-6">
        {/* Welcome Message */}
        <div className="card bg-gradient-to-r from-primary-500 to-primary-600 text-white">
          <h2 className="text-2xl font-bold mb-2">Welcome, {user?.full_name}!</h2>
          <p className="opacity-90">Manage and verify your advisees' achievements here</p>
        </div>

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

        {/* Recent Achievements to Review */}
        <div className="card">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-bold text-gray-900">Recent Achievements</h2>
            <a href="/lecturer/achievements" className="text-primary-600 hover:text-primary-700 font-medium">
              View All
            </a>
          </div>
          {recentAchievements.length > 0 ? (
            <div className="space-y-3">
              {recentAchievements.map((achievement) => (
                <div key={achievement.id} className="p-4 border border-gray-200 rounded-lg hover:border-primary-300 transition-colors">
                  <div className="flex justify-between items-start">
                    <div className="flex-1">
                      <h3 className="font-semibold text-gray-900">{achievement.title}</h3>
                      <p className="text-sm text-gray-600 mt-1">
                        {achievement.student?.name || 'Unknown Student'}
                      </p>
                      <p className="text-xs text-gray-500 mt-1">
                        {achievement.details?.competition_level && (
                          <span className="capitalize">{achievement.details.competition_level} • </span>
                        )}
                        {achievement.achieved_date ? new Date(achievement.achieved_date).toLocaleDateString('id-ID') : 'N/A'}
                      </p>
                    </div>
                    <span className={`badge ${getStatusBadge(achievement.status)}`}>
                      {achievement.status?.replace('_', ' ')}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-12">
              <Award className="w-16 h-16 text-gray-300 mx-auto mb-4" />
              <p className="text-gray-500">No achievements from advisees yet</p>
            </div>
          )}
        </div>

        {/* Quick Actions */}
        <div className="card">
          <h2 className="text-xl font-bold text-gray-900 mb-4">Quick Actions</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <a href="/lecturer/achievements" className="p-4 border-2 border-gray-200 rounded-lg hover:border-primary-500 hover:bg-primary-50 transition-colors">
              <Award className="w-8 h-8 text-primary-600 mb-2" />
              <h3 className="font-semibold text-gray-900">Review Achievements</h3>
              <p className="text-sm text-gray-600">Verify student achievements</p>
              {stats.pending > 0 && (
                <span className="badge badge-warning mt-2">{stats.pending} pending</span>
              )}
            </a>
            <a href="/lecturer/achievements?status=verified" className="p-4 border-2 border-gray-200 rounded-lg hover:border-primary-500 hover:bg-primary-50 transition-colors">
              <CheckCircle className="w-8 h-8 text-green-600 mb-2" />
              <h3 className="font-semibold text-gray-900">Verified Achievements</h3>
              <p className="text-sm text-gray-600">View verified achievements</p>
              <span className="badge badge-success mt-2">{stats.verified} verified</span>
            </a>
          </div>
        </div>
      </div>
    </DashboardLayout>
  )
}
