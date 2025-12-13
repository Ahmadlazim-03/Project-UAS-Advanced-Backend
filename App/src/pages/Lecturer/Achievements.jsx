import { useState, useEffect } from 'react'
import DashboardLayout from '../../components/DashboardLayout'
import { lecturerService, achievementService } from '../../services'
import { CheckCircle, X, FileText, Eye } from 'lucide-react'
import { getStatusBadge, formatDate, getFileUrl } from '../../utils/helpers'

export default function LecturerAchievements() {
  const [achievements, setAchievements] = useState([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState('pending_verification')
  const [selectedAchievement, setSelectedAchievement] = useState(null)
  const [showVerifyModal, setShowVerifyModal] = useState(false)
  const [showDetailModal, setShowDetailModal] = useState(false)
  const [showRejectModal, setShowRejectModal] = useState(false)
  const [comments, setComments] = useState('')
  const [rejectReason, setRejectReason] = useState('')
  const [processing, setProcessing] = useState(false)

  useEffect(() => {
    fetchAchievements()
  }, [filter])

  const fetchAchievements = async () => {
    setLoading(true)
    try {
      // Use the advisee achievements endpoint for lecturers
      const response = await lecturerService.getAdviseeAchievements()
      console.log('Response from backend:', response)
      if (response.status === 'success') {
        // Backend returns paginated response: response.pagination.data.achievements
        const allAchievements = response.pagination?.data?.achievements || []
        console.log('All achievements:', allAchievements)
        
        // Map frontend filter to backend status
        const statusMapping = {
          'draft': 'draft',
          'pending_verification': 'submitted',
          'verified': 'verified',
          'rejected': 'rejected'
        }

        // Filter by status if not 'all'
        let filteredAchievements = allAchievements
        if (filter !== 'all') {
          const backendStatus = statusMapping[filter] || filter
          filteredAchievements = allAchievements.filter(a => a.status === backendStatus)
          console.log(`Filtering by ${filter} (backend: ${backendStatus}), found ${filteredAchievements.length} achievements`)
        }

        setAchievements(filteredAchievements)
      }
    } catch (error) {
      console.error('Error fetching achievements:', error)
    } finally {
      setLoading(false)
    }
  }

  const openVerifyModal = (achievement) => {
    setSelectedAchievement(achievement)
    setComments('')
    setShowVerifyModal(true)
  }

  const openRejectModal = (achievement) => {
    setSelectedAchievement(achievement)
    setRejectReason('')
    setShowRejectModal(true)
  }

  const openDetailModal = (achievement) => {
    setSelectedAchievement(achievement)
    setShowDetailModal(true)
  }

  const getStatusBadge = (status) => {
    const statusMap = {
      draft: { class: 'bg-gray-100 text-gray-800', text: 'Draft' },
      submitted: { class: 'bg-yellow-100 text-yellow-800', text: 'Pending Verification' },
      verified: { class: 'bg-green-100 text-green-800', text: 'Verified' },
      rejected: { class: 'bg-red-100 text-red-800', text: 'Rejected' },
      deleted: { class: 'bg-gray-100 text-gray-500', text: 'Deleted' },
    }
    return statusMap[status] || { class: 'bg-gray-100 text-gray-800', text: status }
  }

  const handleVerify = async () => {
    if (!selectedAchievement) return
    setProcessing(true)

    try {
      // Use mongo_achievement_id for backend API
      const achievementId = selectedAchievement.mongo_achievement_id || selectedAchievement.id
      console.log('Verifying achievement:', achievementId)
      await achievementService.verifyAchievement(achievementId, comments)
      setShowVerifyModal(false)
      setComments('')
      alert('Achievement verified successfully!')
      fetchAchievements()
    } catch (error) {
      console.error('Error verifying achievement:', error)
      alert(error.response?.data?.error || 'Failed to verify achievement')
    } finally {
      setProcessing(false)
    }
  }

  const handleReject = async () => {
    if (!selectedAchievement || !rejectReason.trim()) {
      alert('Please provide a rejection reason')
      return
    }
    setProcessing(true)

    try {
      // Use mongo_achievement_id for backend API
      const achievementId = selectedAchievement.mongo_achievement_id || selectedAchievement.id
      console.log('Rejecting achievement:', achievementId)
      await achievementService.rejectAchievement(achievementId, rejectReason)
      setShowRejectModal(false)
      setRejectReason('')
      alert('Achievement rejected successfully!')
      fetchAchievements()
    } catch (error) {
      console.error('Error rejecting achievement:', error)
      alert(error.response?.data?.error || 'Failed to reject achievement')
    } finally {
      setProcessing(false)
    }
  }

  if (loading) {
    return (
      <DashboardLayout title="Student Achievements">
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
        </div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout title="Student Achievements">
      <div className="space-y-6">
        {/* Filter Tabs */}
        <div className="flex space-x-2 border-b border-gray-200">
          {[
            { value: 'pending_verification', label: 'Pending Review' },
            { value: 'verified', label: 'Verified' },
            { value: 'rejected', label: 'Rejected' },
            { value: 'all', label: 'All' }
          ].map((tab) => (
            <button
              key={tab.value}
              onClick={() => setFilter(tab.value)}
              className={`px-4 py-2 font-medium border-b-2 transition-colors ${filter === tab.value
                  ? 'border-primary-600 text-primary-600'
                  : 'border-transparent text-gray-600 hover:text-gray-900'
                }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* Achievements List */}
        <div className="grid grid-cols-1 gap-6">
          {achievements.length > 0 ? (
            achievements.map((achievement) => {
              const statusBadge = getStatusBadge(achievement.status)
              return (
                <div key={achievement.id} className="card">
                  <div className="flex justify-between items-start mb-4">
                    <div className="flex-1">
                      <div className="flex items-start justify-between mb-2">
                        <h3 className="text-xl font-bold text-gray-900">
                          {achievement.title}
                        </h3>
                        <span className={`badge ${statusBadge.class}`}>
                          {statusBadge.text}
                        </span>
                      </div>
                      <div className="flex items-center space-x-2 mb-2">
                        <span className="text-sm font-medium text-primary-600">
                          {achievement.student?.name || 'Unknown Student'}
                        </span>
                        <span className="text-gray-400">•</span>
                        <span className="text-sm text-gray-500">
                          {achievement.student?.student_id || 'N/A'}
                        </span>
                      </div>
                      <p className="text-gray-600 mb-3">{achievement.description}</p>
                      <div className="flex items-center flex-wrap gap-4 text-sm text-gray-500">
                        <span>📅 {formatDate(achievement.achieved_date)}</span>
                        {achievement.details?.competition_level && (
                          <span className="capitalize">🏆 {achievement.details.competition_level}</span>
                        )}
                        {achievement.details?.rank && <span>🥇 Rank: {achievement.details.rank}</span>}
                        {achievement.details?.medal_type && (
                          <span className="capitalize">🎖️ {achievement.details.medal_type}</span>
                        )}
                      </div>
                    </div>
                  </div>

                  {/* Verification/Rejection Info */}
                  {achievement.verified_at && (
                    <div className="mt-4 p-3 bg-green-50 rounded-lg border border-green-200">
                      <p className="text-sm text-green-800">
                        ✓ Verified on {formatDate(achievement.verified_at)}
                      </p>
                    </div>
                  )}

                  {achievement.rejection_note && (
                    <div className="mt-4 p-3 bg-red-50 rounded-lg border border-red-200">
                      <p className="text-sm text-red-800">
                        ✗ Rejected: {achievement.rejection_note}
                      </p>
                    </div>
                  )}

                  {/* Actions */}
                  <div className="mt-4 flex flex-wrap gap-2">
                    {achievement.status === 'submitted' && (
                      <>
                        <button
                          onClick={() => openVerifyModal(achievement)}
                          className="btn btn-primary flex items-center space-x-2"
                        >
                          <CheckCircle className="w-4 h-4" />
                          <span>Verify</span>
                        </button>
                        <button
                          onClick={() => openRejectModal(achievement)}
                          className="btn btn-danger flex items-center space-x-2"
                        >
                          <X className="w-4 h-4" />
                          <span>Reject</span>
                        </button>
                      </>
                    )}
                    <button
                      onClick={() => openDetailModal(achievement)}
                      className="btn btn-secondary flex items-center space-x-2"
                    >
                      <Eye className="w-4 h-4" />
                      <span>View Details</span>
                    </button>
                  </div>
                </div>
              )
            })
          ) : (
            <div className="card text-center py-12">
              <FileText className="w-16 h-16 text-gray-300 mx-auto mb-4" />
              <p className="text-gray-500">
                {filter === 'pending_verification'
                  ? 'No achievements pending verification'
                  : 'No achievements found'}
              </p>
            </div>
          )}
        </div>
      </div>

      {/* Verification Modal */}
      {showVerifyModal && selectedAchievement && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg max-w-2xl w-full">
            <div className="flex justify-between items-center p-6 border-b">
              <h2 className="text-2xl font-bold">Verify Achievement</h2>
              <button
                onClick={() => setShowVerifyModal(false)}
                className="text-gray-400 hover:text-gray-600"
              >
                <X className="w-6 h-6" />
              </button>
            </div>
            <div className="p-6">
              <div className="mb-4 p-4 bg-gray-50 rounded-lg">
                <h3 className="font-semibold text-gray-900 mb-2">{selectedAchievement.title}</h3>
                <p className="text-sm text-gray-600 mb-2">
                  Student: {selectedAchievement.student?.user?.full_name}
                </p>
                <p className="text-gray-600">{selectedAchievement.description}</p>
              </div>
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Comments (Optional)
                </label>
                <textarea
                  value={comments}
                  onChange={(e) => setComments(e.target.value)}
                  className="input"
                  rows="4"
                  placeholder="Add verification comments..."
                />
              </div>
              <div className="flex justify-end space-x-3">
                <button
                  onClick={() => setShowVerifyModal(false)}
                  className="btn btn-secondary"
                  disabled={processing}
                >
                  Cancel
                </button>
                <button
                  onClick={handleVerify}
                  className="btn btn-primary"
                  disabled={processing}
                >
                  {processing ? 'Verifying...' : 'Confirm Verification'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Rejection Modal */}
      {showRejectModal && selectedAchievement && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg max-w-2xl w-full">
            <div className="flex justify-between items-center p-6 border-b">
              <h2 className="text-2xl font-bold text-red-600">Reject Achievement</h2>
              <button
                onClick={() => setShowRejectModal(false)}
                className="text-gray-400 hover:text-gray-600"
              >
                <X className="w-6 h-6" />
              </button>
            </div>
            <div className="p-6">
              <div className="mb-4 p-4 bg-gray-50 rounded-lg">
                <h3 className="font-semibold text-gray-900 mb-2">{selectedAchievement.title}</h3>
                <p className="text-sm text-gray-600 mb-2">
                  Student: {selectedAchievement.student?.user?.full_name}
                </p>
              </div>
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Rejection Reason *
                </label>
                <textarea
                  value={rejectReason}
                  onChange={(e) => setRejectReason(e.target.value)}
                  className="input"
                  rows="4"
                  placeholder="Please provide a detailed reason for rejection..."
                  required
                />
                <p className="text-xs text-gray-500 mt-1">
                  This will be shown to the student so they can address the issues.
                </p>
              </div>
              <div className="flex justify-end space-x-3">
                <button
                  onClick={() => setShowRejectModal(false)}
                  className="btn btn-secondary"
                  disabled={processing}
                >
                  Cancel
                </button>
                <button
                  onClick={handleReject}
                  className="btn btn-danger"
                  disabled={processing || !rejectReason.trim()}
                >
                  {processing ? 'Rejecting...' : 'Confirm Rejection'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* View Details Modal */}
      {showDetailModal && selectedAchievement && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg max-w-2xl w-full max-h-[90vh] overflow-y-auto">
            <div className="flex justify-between items-center p-6 border-b sticky top-0 bg-white">
              <h2 className="text-2xl font-bold">Achievement Details</h2>
              <button
                onClick={() => setShowDetailModal(false)}
                className="text-gray-400 hover:text-gray-600"
              >
                <X className="w-6 h-6" />
              </button>
            </div>

            <div className="p-6 space-y-6">
              {/* Header */}
              <div>
                <div className="flex justify-between items-start mb-2">
                  <h3 className="text-xl font-bold text-gray-900">{selectedAchievement.title}</h3>
                  <span className={`badge ${getStatusBadge(selectedAchievement.status).class}`}>
                    {getStatusBadge(selectedAchievement.status).text}
                  </span>
                </div>
                <p className="text-gray-600">{selectedAchievement.description}</p>
              </div>

              {/* Student Info */}
              <div className="p-4 bg-blue-50 rounded-lg border border-blue-200">
                <h4 className="font-medium text-blue-800 mb-2">Student Information</h4>
                <p className="text-sm text-blue-700">
                  Name: {selectedAchievement.student?.user?.full_name || 'Unknown'}
                </p>
                <p className="text-sm text-blue-700">
                  Student ID: {selectedAchievement.student?.student_id}
                </p>
                <p className="text-sm text-blue-700">
                  Program: {selectedAchievement.student?.program_study}
                </p>
              </div>

              {/* Info Grid */}
              <div className="grid grid-cols-2 gap-4">
                <div className="p-3 bg-gray-50 rounded-lg">
                  <p className="text-xs text-gray-500">Achievement Date</p>
                  <p className="font-medium">{formatDate(selectedAchievement.achieved_date)}</p>
                </div>
                {selectedAchievement.details?.competition_level && (
                  <div className="p-3 bg-gray-50 rounded-lg">
                    <p className="text-xs text-gray-500">Level</p>
                    <p className="font-medium capitalize">{selectedAchievement.details.competition_level}</p>
                  </div>
                )}
                {selectedAchievement.details?.competition_name && (
                  <div className="p-3 bg-gray-50 rounded-lg col-span-2">
                    <p className="text-xs text-gray-500">Competition Name</p>
                    <p className="font-medium">{selectedAchievement.details.competition_name}</p>
                  </div>
                )}
                {selectedAchievement.details?.rank && (
                  <div className="p-3 bg-gray-50 rounded-lg">
                    <p className="text-xs text-gray-500">Rank</p>
                    <p className="font-medium">{selectedAchievement.details.rank}</p>
                  </div>
                )}
                {selectedAchievement.details?.medal_type && (
                  <div className="p-3 bg-gray-50 rounded-lg">
                    <p className="text-xs text-gray-500">Medal</p>
                    <p className="font-medium capitalize">{selectedAchievement.details.medal_type}</p>
                  </div>
                )}
                {selectedAchievement.details?.organizer && (
                  <div className="p-3 bg-gray-50 rounded-lg">
                    <p className="text-xs text-gray-500">Organizer</p>
                    <p className="font-medium">{selectedAchievement.details.organizer}</p>
                  </div>
                )}
                {selectedAchievement.details?.location && (
                  <div className="p-3 bg-gray-50 rounded-lg">
                    <p className="text-xs text-gray-500">Location</p>
                    <p className="font-medium">{selectedAchievement.details.location}</p>
                  </div>
                )}
              </div>

              {/* Certificate Link */}
              {selectedAchievement.attachments && selectedAchievement.attachments.length > 0 && (
                <div className="p-4 bg-purple-50 rounded-lg border border-purple-200">
                  <h4 className="font-medium text-purple-800 mb-2">Attachments</h4>
                  {selectedAchievement.attachments.map((attachment, index) => (
                    <a
                      key={index}
                      href={getFileUrl(attachment.url)}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-sm text-purple-600 hover:underline block"
                    >
                      {attachment.filename || `Attachment ${index + 1}`} →
                    </a>
                  ))}
                </div>
              )}

              {/* Verification Info */}
              {selectedAchievement.verified_at && (
                <div className="p-4 bg-green-50 rounded-lg border border-green-200">
                  <h4 className="font-medium text-green-800 mb-2">Verification Details</h4>
                  <p className="text-sm text-green-700">
                    Verified on: {formatDate(selectedAchievement.verified_at)}
                  </p>
                </div>
              )}

              {selectedAchievement.rejection_note && (
                <div className="p-4 bg-red-50 rounded-lg border border-red-200">
                  <h4 className="font-medium text-red-800 mb-2">Rejection Details</h4>
                  <p className="text-sm text-red-700">
                    Reason: {selectedAchievement.rejection_note}
                  </p>
                </div>
              )}

              {/* Timestamps */}
              <div className="border-t pt-4 text-sm text-gray-500">
                <p>Created: {formatDate(selectedAchievement.created_at)}</p>
                <p>Last Updated: {formatDate(selectedAchievement.updated_at)}</p>
                {selectedAchievement.submitted_at && (
                  <p>Submitted: {formatDate(selectedAchievement.submitted_at)}</p>
                )}
              </div>

              {/* Actions for pending achievements */}
              {selectedAchievement.status === 'submitted' && (
                <div className="flex justify-end space-x-3 pt-4 border-t">
                  <button
                    onClick={() => {
                      setShowDetailModal(false)
                      openRejectModal(selectedAchievement)
                    }}
                    className="btn btn-danger"
                  >
                    Reject
                  </button>
                  <button
                    onClick={() => {
                      setShowDetailModal(false)
                      openVerifyModal(selectedAchievement)
                    }}
                    className="btn btn-primary"
                  >
                    Verify
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </DashboardLayout>
  )
}
