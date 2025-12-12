import { useState, useEffect } from 'react'
import DashboardLayout from '../../components/DashboardLayout'
import { achievementService } from '../../services'
import { Award, CheckCircle, XCircle, Eye, Calendar, Trophy, Medal } from 'lucide-react'

export default function AdminAchievements() {
  const [achievements, setAchievements] = useState([])
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState('submitted')
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
    totalPages: 0,
  })
  const [selectedAchievement, setSelectedAchievement] = useState(null)
  const [showVerifyModal, setShowVerifyModal] = useState(false)
  const [showRejectModal, setShowRejectModal] = useState(false)
  const [comments, setComments] = useState('')
  const [reason, setReason] = useState('')

  useEffect(() => {
    fetchAchievements()
  }, [activeTab, pagination.page])

  const fetchAchievements = async () => {
    try {
      setLoading(true)
      const statusMap = {
        submitted: 'submitted',
        verified: 'verified',
        rejected: 'rejected',
        all: '',
      }
      
      const params = {
        page: pagination.page,
        limit: pagination.limit,
        status: statusMap[activeTab],
      }

      const response = await achievementService.getAchievements(params)
      console.log('Admin Achievements Response:', response)
      
      const achievementsList = response.pagination?.data?.achievements || []
      setAchievements(achievementsList)
      
      if (response.pagination) {
        setPagination({
          page: response.pagination.page,
          limit: response.pagination.limit,
          total: response.pagination.total,
          totalPages: response.pagination.totalPages,
        })
      }
    } catch (error) {
      console.error('Error fetching achievements:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleVerify = async () => {
    try {
      await achievementService.verifyAchievement(
        selectedAchievement.mongo_achievement_id,
        comments
      )
      setShowVerifyModal(false)
      setComments('')
      fetchAchievements()
    } catch (error) {
      console.error('Error verifying achievement:', error)
      alert('Failed to verify achievement')
    }
  }

  const handleReject = async () => {
    try {
      await achievementService.rejectAchievement(
        selectedAchievement.mongo_achievement_id,
        reason
      )
      setShowRejectModal(false)
      setReason('')
      fetchAchievements()
    } catch (error) {
      console.error('Error rejecting achievement:', error)
      alert('Failed to reject achievement')
    }
  }

  const tabs = [
    { key: 'submitted', label: 'Pending Review', count: achievements.length },
    { key: 'verified', label: 'Verified', count: achievements.length },
    { key: 'rejected', label: 'Rejected', count: achievements.length },
    { key: 'all', label: 'All', count: achievements.length },
  ]

  if (loading) {
    return (
      <DashboardLayout title="Achievement Management">
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
        </div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout title="Achievement Management">
      <div className="space-y-6">
        {/* Tabs */}
        <div className="card">
          <div className="border-b border-gray-200">
            <nav className="-mb-px flex space-x-8">
              {tabs.map((tab) => (
                <button
                  key={tab.key}
                  onClick={() => {
                    setActiveTab(tab.key)
                    setPagination((prev) => ({ ...prev, page: 1 }))
                  }}
                  className={`
                    py-4 px-1 border-b-2 font-medium text-sm
                    ${
                      activeTab === tab.key
                        ? 'border-primary-500 text-primary-600'
                        : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                    }
                  `}
                >
                  {tab.label}
                  <span className="ml-2 py-0.5 px-2 rounded-full text-xs bg-gray-100">
                    {tab.count}
                  </span>
                </button>
              ))}
            </nav>
          </div>
        </div>

        {/* Achievements List */}
        <div className="grid grid-cols-1 gap-6">
          {achievements.length === 0 ? (
            <div className="card text-center py-12">
              <Award className="w-16 h-16 text-gray-400 mx-auto mb-4" />
              <p className="text-gray-600">No achievements found</p>
            </div>
          ) : (
            achievements.map((achievement) => (
              <div key={achievement.id} className="card hover:shadow-lg transition-shadow">
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-3 mb-2">
                      <Trophy className="w-5 h-5 text-primary-600" />
                      <h3 className="text-lg font-semibold text-gray-900">
                        {achievement.title}
                      </h3>
                      <span
                        className={`px-3 py-1 rounded-full text-xs font-medium ${
                          achievement.status === 'verified'
                            ? 'bg-green-100 text-green-800'
                            : achievement.status === 'rejected'
                            ? 'bg-red-100 text-red-800'
                            : achievement.status === 'submitted'
                            ? 'bg-yellow-100 text-yellow-800'
                            : 'bg-gray-100 text-gray-800'
                        }`}
                      >
                        {achievement.status === 'submitted' ? 'Pending Review' : achievement.status}
                      </span>
                    </div>

                    {/* Student Info */}
                    <div className="mb-3 p-3 bg-gray-50 rounded-lg">
                      <div className="grid grid-cols-2 gap-2 text-sm">
                        <div>
                          <span className="text-gray-600">Student:</span>
                          <span className="ml-2 font-medium text-gray-900">
                            {achievement.student?.name || 'Unknown Student'}
                          </span>
                        </div>
                        <div>
                          <span className="text-gray-600">NIM:</span>
                          <span className="ml-2 font-medium text-gray-900">
                            {achievement.student?.student_id || 'N/A'}
                          </span>
                        </div>
                        <div>
                          <span className="text-gray-600">Email:</span>
                          <span className="ml-2 font-medium text-gray-900">
                            {achievement.student?.email || 'N/A'}
                          </span>
                        </div>
                        <div>
                          <span className="text-gray-600">Program:</span>
                          <span className="ml-2 font-medium text-gray-900">
                            {achievement.student?.program || 'N/A'}
                          </span>
                        </div>
                      </div>
                    </div>

                    {/* Achievement Details */}
                    <p className="text-gray-600 mb-3">{achievement.description}</p>
                    
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-3 text-sm">
                      <div className="flex items-center gap-2">
                        <Medal className="w-4 h-4 text-gray-500" />
                        <span className="text-gray-600">
                          {achievement.details?.competition_level || 'N/A'}
                        </span>
                      </div>
                      <div className="flex items-center gap-2">
                        <Trophy className="w-4 h-4 text-gray-500" />
                        <span className="text-gray-600">
                          Rank: {achievement.details?.rank || 'N/A'}
                        </span>
                      </div>
                      <div className="flex items-center gap-2">
                        <Award className="w-4 h-4 text-gray-500" />
                        <span className="text-gray-600">
                          {achievement.details?.medal_type || 'N/A'}
                        </span>
                      </div>
                      <div className="flex items-center gap-2">
                        <Calendar className="w-4 h-4 text-gray-500" />
                        <span className="text-gray-600">
                          {achievement.achieved_date
                            ? new Date(achievement.achieved_date).toLocaleDateString('id-ID')
                            : 'N/A'}
                        </span>
                      </div>
                    </div>

                    {/* Verified/Rejected Info */}
                    {achievement.status === 'verified' && achievement.verified_by && (
                      <div className="mt-3 p-2 bg-green-50 rounded text-sm text-green-800">
                        Verified by {achievement.verified_by} on{' '}
                        {new Date(achievement.verified_at).toLocaleDateString('id-ID')}
                        {achievement.comments && ` - ${achievement.comments}`}
                      </div>
                    )}
                    {achievement.status === 'rejected' && achievement.rejected_by && (
                      <div className="mt-3 p-2 bg-red-50 rounded text-sm text-red-800">
                        Rejected by {achievement.rejected_by} on{' '}
                        {new Date(achievement.rejected_at).toLocaleDateString('id-ID')}
                        {achievement.rejection_reason && ` - ${achievement.rejection_reason}`}
                      </div>
                    )}
                  </div>

                  {/* Action Buttons */}
                  {achievement.status === 'submitted' && (
                    <div className="flex gap-2 ml-4">
                      <button
                        onClick={() => {
                          setSelectedAchievement(achievement)
                          setShowVerifyModal(true)
                        }}
                        className="btn-primary flex items-center gap-2"
                      >
                        <CheckCircle className="w-4 h-4" />
                        Verify
                      </button>
                      <button
                        onClick={() => {
                          setSelectedAchievement(achievement)
                          setShowRejectModal(true)
                        }}
                        className="btn-danger flex items-center gap-2"
                      >
                        <XCircle className="w-4 h-4" />
                        Reject
                      </button>
                    </div>
                  )}
                </div>
              </div>
            ))
          )}
        </div>

        {/* Pagination */}
        {pagination.totalPages > 1 && (
          <div className="flex justify-center gap-2">
            <button
              onClick={() => setPagination((prev) => ({ ...prev, page: prev.page - 1 }))}
              disabled={pagination.page === 1}
              className="btn-secondary disabled:opacity-50"
            >
              Previous
            </button>
            <span className="px-4 py-2">
              Page {pagination.page} of {pagination.totalPages}
            </span>
            <button
              onClick={() => setPagination((prev) => ({ ...prev, page: prev.page + 1 }))}
              disabled={pagination.page >= pagination.totalPages}
              className="btn-secondary disabled:opacity-50"
            >
              Next
            </button>
          </div>
        )}
      </div>

      {/* Verify Modal */}
      {showVerifyModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-md w-full">
            <h3 className="text-xl font-bold text-gray-900 mb-4">Verify Achievement</h3>
            <p className="text-gray-600 mb-4">
              Are you sure you want to verify "{selectedAchievement?.title}"?
            </p>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Comments (Optional)
              </label>
              <textarea
                value={comments}
                onChange={(e) => setComments(e.target.value)}
                className="input-field"
                rows="3"
                placeholder="Add any comments..."
              />
            </div>
            <div className="flex gap-3">
              <button onClick={handleVerify} className="btn-primary flex-1">
                Verify
              </button>
              <button
                onClick={() => {
                  setShowVerifyModal(false)
                  setComments('')
                }}
                className="btn-secondary flex-1"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Reject Modal */}
      {showRejectModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-md w-full">
            <h3 className="text-xl font-bold text-gray-900 mb-4">Reject Achievement</h3>
            <p className="text-gray-600 mb-4">
              Please provide a reason for rejecting "{selectedAchievement?.title}":
            </p>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Reason <span className="text-red-500">*</span>
              </label>
              <textarea
                value={reason}
                onChange={(e) => setReason(e.target.value)}
                className="input-field"
                rows="3"
                placeholder="Provide a detailed reason..."
                required
              />
            </div>
            <div className="flex gap-3">
              <button
                onClick={handleReject}
                disabled={!reason.trim()}
                className="btn-danger flex-1 disabled:opacity-50"
              >
                Reject
              </button>
              <button
                onClick={() => {
                  setShowRejectModal(false)
                  setReason('')
                }}
                className="btn-secondary flex-1"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </DashboardLayout>
  )
}
