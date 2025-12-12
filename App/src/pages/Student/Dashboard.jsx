import { useState, useEffect } from 'react'
import DashboardLayout from '../../components/DashboardLayout'
import { achievementService } from '../../services'
import { Award, TrendingUp, Clock, CheckCircle } from 'lucide-react'
import { useAuth } from '../../contexts/AuthContext'

export default function StudentDashboard() {
  const { user } = useAuth()
  const [stats, setStats] = useState({
    total: 0,
    draft: 0,
    pending: 0,
    verified: 0,
  })
  const [recentAchievements, setRecentAchievements] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchAchievements()
  }, [])

  const fetchAchievements = async () => {
    try {
      const response = await achievementService.getAchievements({ page: 1, limit: 5 })
      if (response.status === 'success' && response.data) {
        const achievements = response.data.achievements || []
        setRecentAchievements(achievements)
        
        // Calculate stats
        const statsData = {
          total: achievements.length,
          draft: achievements.filter(a => a.status === 'draft').length,
          pending: achievements.filter(a => a.status === 'pending_verification').length,
          verified: achievements.filter(a => a.status === 'verified').length,
        }
        setStats(statsData)
      }
    } catch (error) {
      console.error('Error fetching achievements:', error)
    } finally {
      setLoading(false)
    }
  }

  const statCards = [
    {
      title: 'Total Achievements',
      value: stats.total,
      icon: Award,
      color: 'bg-purple-500',
    },
    {
      title: 'Verified',
      value: stats.verified,
      icon: CheckCircle,
      color: 'bg-green-500',
    },
    {
      title: 'Pending',
      value: stats.pending,
      icon: Clock,
      color: 'bg-yellow-500',
    },
    {
      title: 'Draft',
      value: stats.draft,
      icon: TrendingUp,
      color: 'bg-blue-500',
    },
  ]

  if (loading) {
    return (
      <DashboardLayout title="Student Dashboard">
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
        </div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout title="Student Dashboard">
      <div className="space-y-6">
        {/* Welcome Message */}
        <div className="card bg-gradient-to-r from-primary-500 to-primary-600 text-white">
          <h2 className="text-2xl font-bold mb-2">Welcome back, {user?.full_name}!</h2>
          <p className="opacity-90">Track and manage your achievements here</p>
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

        {/* Recent Achievements */}
        <div className="card">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-bold text-gray-900">Recent Achievements</h2>
            <a href="/student/achievements" className="text-primary-600 hover:text-primary-700 font-medium">
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
                      <p className="text-sm text-gray-600 mt-1">{achievement.description}</p>
                    </div>
                    <span className={`badge ${
                      achievement.status === 'verified' ? 'badge-success' :
                      achievement.status === 'pending_verification' ? 'badge-warning' :
                      achievement.status === 'rejected' ? 'badge-danger' :
                      'badge-info'
                    }`}>
                      {achievement.status?.replace('_', ' ')}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-12">
              <Award className="w-16 h-16 text-gray-300 mx-auto mb-4" />
              <p className="text-gray-500">No achievements yet</p>
              <a href="/student/achievements" className="btn btn-primary mt-4 inline-block">
                Add Your First Achievement
              </a>
            </div>
          )}
        </div>
      </div>
    </DashboardLayout>
  )
}
