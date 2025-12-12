import { useState, useEffect } from 'react'
import DashboardLayout from '../../components/DashboardLayout'
import { achievementService, fileService } from '../../services'
import { Plus, Edit, Trash2, Send, FileText, X, Upload, Eye } from 'lucide-react'
import { getStatusBadge, formatDate, getFileUrl } from '../../utils/helpers'

export default function StudentAchievements() {
  const [achievements, setAchievements] = useState([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [showDetailModal, setShowDetailModal] = useState(false)
  const [selectedAchievement, setSelectedAchievement] = useState(null)
  const [editingAchievement, setEditingAchievement] = useState(null)
  const [saving, setSaving] = useState(false)
  const [uploadingFile, setUploadingFile] = useState(false)

  const [formData, setFormData] = useState({
    title: '',
    description: '',
    achieved_date: '',
    data: {
      type: 'competition', // competition, publication, research_grant, certification
      // Competition fields
      competition_name: '',
      competition_level: 'national', // international, national, regional, university
      rank: '',
      medal_type: '', // gold, silver, bronze
      organizer: '',
      location: '',
      participants_count: '',
      // Publication fields
      publication_type: '', // journal, conference, book
      journal_name: '',
      publisher: '',
      doi: '',
      issn: '',
      volume: '',
      issue: '',
      pages: '',
      // Common
      certificate_url: '',
    }
  })

  useEffect(() => {
    fetchAchievements()
  }, [])

  const fetchAchievements = async () => {
    try {
      const response = await achievementService.getAchievements()
      console.log('API Response:', response)
      
      // Response structure: { status: 'success', pagination: { data: { achievements: [...] } } }
      if (response && response.pagination && response.pagination.data && response.pagination.data.achievements) {
        console.log('Achievements Data:', response.pagination.data.achievements)
        setAchievements(response.pagination.data.achievements)
      } else {
        console.log('No achievements found or invalid response structure')
        setAchievements([])
      }
    } catch (error) {
      console.error('Error fetching achievements:', error)
      setAchievements([])
    } finally {
      setLoading(false)
    }
  }

  const resetForm = () => {
    setFormData({
      title: '',
      description: '',
      achieved_date: '',
      data: {
        type: 'competition',
        competition_name: '',
        competition_level: 'national',
        rank: '',
        medal_type: '',
        organizer: '',
        location: '',
        participants_count: '',
        publication_type: '',
        journal_name: '',
        publisher: '',
        doi: '',
        issn: '',
        volume: '',
        issue: '',
        pages: '',
        certificate_url: '',
      }
    })
  }

  const openCreateModal = () => {
    setEditingAchievement(null)
    resetForm()
    setShowModal(true)
  }

  const openEditModal = (achievement) => {
    setEditingAchievement({ ...achievement, id: achievement.mongo_achievement_id })
    
    // Get details from backend (it's in details field)
    const details = achievement.details || {}
    
    setFormData({
      title: achievement.title || '',
      description: achievement.description || '',
      achieved_date: achievement.achieved_date ? achievement.achieved_date.split('T')[0] : '',
      data: {
        type: details.competition_name ? 'competition' : details.publication_type ? 'publication' : 'competition',
        competition_name: details.competition_name || '',
        competition_level: details.competition_level || 'national',
        rank: details.rank || '',
        medal_type: details.medal_type || '',
        organizer: details.organizer || '',
        location: details.location || '',
        participants_count: details.participants_count || '',
        publication_type: details.publication_type || '',
        journal_name: details.journal_name || '',
        publisher: details.publisher || '',
        doi: details.doi || '',
        issn: details.issn || '',
        volume: details.volume || '',
        issue: details.issue || '',
        pages: details.pages || '',
        certificate_url: details.certificate_url || '',
      }
    })
    setShowModal(true)
  }

  const openDetailModal = (achievement) => {
    setSelectedAchievement(achievement)
    setShowDetailModal(true)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setSaving(true)

    try {
      // Build data object based on type
      const achievementData = {
        title: formData.title,
        description: formData.description,
        achieved_date: formData.achieved_date,
        data: {}
      }

      if (formData.data.type === 'competition') {
        achievementData.data = {
          competition_name: formData.data.competition_name,
          competition_level: formData.data.competition_level,
          rank: formData.data.rank ? parseInt(formData.data.rank) : undefined,
          medal_type: formData.data.medal_type || undefined,
          organizer: formData.data.organizer || undefined,
          location: formData.data.location || undefined,
          participants_count: formData.data.participants_count ? parseInt(formData.data.participants_count) : undefined,
          certificate_url: formData.data.certificate_url || undefined,
        }
      } else if (formData.data.type === 'publication') {
        achievementData.data = {
          publication_type: formData.data.publication_type,
          journal_name: formData.data.journal_name,
          publisher: formData.data.publisher || undefined,
          doi: formData.data.doi || undefined,
          issn: formData.data.issn || undefined,
          volume: formData.data.volume || undefined,
          issue: formData.data.issue || undefined,
          pages: formData.data.pages || undefined,
          certificate_url: formData.data.certificate_url || undefined,
        }
      }

      // Remove undefined values
      Object.keys(achievementData.data).forEach(key => {
        if (achievementData.data[key] === undefined) {
          delete achievementData.data[key]
        }
      })

      if (editingAchievement) {
        await achievementService.updateAchievement(editingAchievement.id, achievementData)
      } else {
        await achievementService.createAchievement(achievementData)
      }

      setShowModal(false)
      resetForm()
      fetchAchievements()
    } catch (error) {
      console.error('Error saving achievement:', error)
      alert(error.response?.data?.message || 'Failed to save achievement')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id) => {
    if (!confirm('Are you sure you want to delete this achievement?')) return
    try {
      await achievementService.deleteAchievement(id)
      fetchAchievements()
    } catch (error) {
      console.error('Error deleting achievement:', error)
      alert('Failed to delete achievement')
    }
  }

  const handleSubmitForVerification = async (id) => {
    if (!confirm('Submit this achievement for verification?')) return
    try {
      await achievementService.submitAchievement(id)
      fetchAchievements()
    } catch (error) {
      console.error('Error submitting achievement:', error)
      alert(error.response?.data?.message || 'Failed to submit achievement')
    }
  }

  const handleFileUpload = async (e) => {
    const file = e.target.files[0]
    if (!file) return

    setUploadingFile(true)
    try {
      const response = await fileService.uploadFile(file)
      if (response.status === 'success') {
        setFormData({
          ...formData,
          data: {
            ...formData.data,
            certificate_url: response.data.url || response.data.filename
          }
        })
      }
    } catch (error) {
      console.error('Error uploading file:', error)
      alert('Failed to upload file')
    } finally {
      setUploadingFile(false)
    }
  }

  if (loading) {
    return (
      <DashboardLayout title="My Achievements">
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
        </div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout title="My Achievements">
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <p className="text-gray-600">Manage your academic achievements</p>
          <button
            onClick={openCreateModal}
            className="btn btn-primary flex items-center space-x-2"
          >
            <Plus className="w-5 h-5" />
            <span>Add Achievement</span>
          </button>
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
                      <h3 className="text-xl font-bold text-gray-900 mb-2">
                        {achievement.title}
                      </h3>
                      <p className="text-gray-600 mb-3">{achievement.description}</p>
                      <div className="flex items-center space-x-4 text-sm text-gray-500">
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
                    <span className={`badge ${statusBadge.class}`}>
                      {statusBadge.text}
                    </span>
                  </div>

                  {/* Verification Info */}
                  {achievement.verification && (
                    <div className="mt-4 p-3 bg-green-50 rounded-lg border border-green-200">
                      <p className="text-sm text-green-800">
                        ✓ Verified {achievement.verification.comments && `: ${achievement.verification.comments}`}
                      </p>
                    </div>
                  )}

                  {achievement.rejection && (
                    <div className="mt-4 p-3 bg-red-50 rounded-lg border border-red-200">
                      <p className="text-sm text-red-800">
                        ✗ Rejected: {achievement.rejection.reason}
                      </p>
                    </div>
                  )}

                  {/* Actions */}
                  <div className="mt-4 flex flex-wrap gap-2">
                    {achievement.status === 'draft' && (
                      <>
                        <button
                          onClick={() => handleSubmitForVerification(achievement.mongo_achievement_id)}
                          className="btn btn-primary flex items-center space-x-2"
                        >
                          <Send className="w-4 h-4" />
                          <span>Submit for Verification</span>
                        </button>
                        <button
                          onClick={() => openEditModal(achievement)}
                          className="btn btn-secondary flex items-center space-x-2"
                        >
                          <Edit className="w-4 h-4" />
                          <span>Edit</span>
                        </button>
                        <button
                          onClick={() => handleDelete(achievement.mongo_achievement_id)}
                          className="btn btn-danger flex items-center space-x-2"
                        >
                          <Trash2 className="w-4 h-4" />
                          <span>Delete</span>
                        </button>
                      </>
                    )}
                    {achievement.status === 'rejected' && (
                      <button
                        onClick={() => openEditModal(achievement)}
                        className="btn btn-primary flex items-center space-x-2"
                      >
                        <Edit className="w-4 h-4" />
                        <span>Revise & Resubmit</span>
                      </button>
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
              <p className="text-gray-500 mb-4">No achievements yet</p>
              <button
                onClick={openCreateModal}
                className="btn btn-primary inline-flex items-center space-x-2"
              >
                <Plus className="w-5 h-5" />
                <span>Add Your First Achievement</span>
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Create/Edit Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg max-w-3xl w-full max-h-[90vh] overflow-y-auto">
            <div className="flex justify-between items-center p-6 border-b sticky top-0 bg-white">
              <h2 className="text-2xl font-bold">
                {editingAchievement ? 'Edit Achievement' : 'Add New Achievement'}
              </h2>
              <button
                onClick={() => setShowModal(false)}
                className="text-gray-400 hover:text-gray-600"
              >
                <X className="w-6 h-6" />
              </button>
            </div>

            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              {/* Achievement Type */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Achievement Type *
                </label>
                <select
                  value={formData.data.type}
                  onChange={(e) => setFormData({
                    ...formData,
                    data: { ...formData.data, type: e.target.value }
                  })}
                  className="input"
                >
                  <option value="competition">Competition / Lomba</option>
                  <option value="publication">Publication / Publikasi</option>
                  <option value="research_grant">Research Grant / Hibah</option>
                  <option value="certification">Certification / Sertifikasi</option>
                </select>
              </div>

              {/* Basic Info */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Title *
                </label>
                <input
                  type="text"
                  required
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  className="input"
                  placeholder="Achievement title"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Description *
                </label>
                <textarea
                  required
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  className="input"
                  rows="3"
                  placeholder="Describe your achievement"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Achievement Date *
                </label>
                <input
                  type="date"
                  required
                  value={formData.achieved_date}
                  onChange={(e) => setFormData({ ...formData, achieved_date: e.target.value })}
                  className="input"
                />
              </div>

              {/* Competition Fields */}
              {formData.data.type === 'competition' && (
                <div className="border-t pt-4 space-y-4">
                  <h3 className="text-lg font-medium text-gray-900">Competition Details</h3>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Competition Name
                      </label>
                      <input
                        type="text"
                        value={formData.data.competition_name}
                        onChange={(e) => setFormData({
                          ...formData,
                          data: { ...formData.data, competition_name: e.target.value }
                        })}
                        className="input"
                        placeholder="e.g., National Coding Championship"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Competition Level
                      </label>
                      <select
                        value={formData.data.competition_level}
                        onChange={(e) => setFormData({
                          ...formData,
                          data: { ...formData.data, competition_level: e.target.value }
                        })}
                        className="input"
                      >
                        <option value="international">International</option>
                        <option value="national">National</option>
                        <option value="regional">Regional</option>
                        <option value="university">University</option>
                      </select>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Rank / Position
                      </label>
                      <input
                        type="number"
                        value={formData.data.rank}
                        onChange={(e) => setFormData({
                          ...formData,
                          data: { ...formData.data, rank: e.target.value }
                        })}
                        className="input"
                        placeholder="e.g., 1"
                        min="1"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Medal Type
                      </label>
                      <select
                        value={formData.data.medal_type}
                        onChange={(e) => setFormData({
                          ...formData,
                          data: { ...formData.data, medal_type: e.target.value }
                        })}
                        className="input"
                      >
                        <option value="">Select Medal</option>
                        <option value="gold">Gold</option>
                        <option value="silver">Silver</option>
                        <option value="bronze">Bronze</option>
                      </select>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Organizer
                      </label>
                      <input
                        type="text"
                        value={formData.data.organizer}
                        onChange={(e) => setFormData({
                          ...formData,
                          data: { ...formData.data, organizer: e.target.value }
                        })}
                        className="input"
                        placeholder="e.g., Ministry of Education"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Location
                      </label>
                      <input
                        type="text"
                        value={formData.data.location}
                        onChange={(e) => setFormData({
                          ...formData,
                          data: { ...formData.data, location: e.target.value }
                        })}
                        className="input"
                        placeholder="e.g., Jakarta, Indonesia"
                      />
                    </div>
                  </div>
                </div>
              )}

              {/* Publication Fields */}
              {formData.data.type === 'publication' && (
                <div className="border-t pt-4 space-y-4">
                  <h3 className="text-lg font-medium text-gray-900">Publication Details</h3>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Publication Type
                      </label>
                      <select
                        value={formData.data.publication_type}
                        onChange={(e) => setFormData({
                          ...formData,
                          data: { ...formData.data, publication_type: e.target.value }
                        })}
                        className="input"
                      >
                        <option value="">Select Type</option>
                        <option value="journal">Journal Article</option>
                        <option value="conference">Conference Paper</option>
                        <option value="book">Book / Book Chapter</option>
                      </select>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Journal/Conference Name
                      </label>
                      <input
                        type="text"
                        value={formData.data.journal_name}
                        onChange={(e) => setFormData({
                          ...formData,
                          data: { ...formData.data, journal_name: e.target.value }
                        })}
                        className="input"
                        placeholder="e.g., IEEE Transactions"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        DOI
                      </label>
                      <input
                        type="text"
                        value={formData.data.doi}
                        onChange={(e) => setFormData({
                          ...formData,
                          data: { ...formData.data, doi: e.target.value }
                        })}
                        className="input"
                        placeholder="e.g., 10.1109/example.2024.123456"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Publisher
                      </label>
                      <input
                        type="text"
                        value={formData.data.publisher}
                        onChange={(e) => setFormData({
                          ...formData,
                          data: { ...formData.data, publisher: e.target.value }
                        })}
                        className="input"
                        placeholder="e.g., IEEE"
                      />
                    </div>
                  </div>
                </div>
              )}

              {/* Certificate Upload */}
              <div className="border-t pt-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Certificate / Supporting Document
                </label>
                <div className="flex items-center space-x-4">
                  <label className="btn btn-secondary cursor-pointer">
                    <Upload className="w-4 h-4 mr-2" />
                    {uploadingFile ? 'Uploading...' : 'Upload File'}
                    <input
                      type="file"
                      className="hidden"
                      onChange={handleFileUpload}
                      accept=".pdf,.jpg,.jpeg,.png"
                      disabled={uploadingFile}
                    />
                  </label>
                  {formData.data.certificate_url && (
                    <span className="text-sm text-green-600">
                      ✓ File uploaded: {formData.data.certificate_url}
                    </span>
                  )}
                </div>
                <p className="text-xs text-gray-500 mt-1">
                  Accepted formats: PDF, JPG, PNG (max 5MB)
                </p>
              </div>

              <div className="flex justify-end space-x-3 pt-4 border-t">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="btn btn-secondary"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="btn btn-primary"
                  disabled={saving}
                >
                  {saving ? 'Saving...' : (editingAchievement ? 'Update Achievement' : 'Create Achievement')}
                </button>
              </div>
            </form>
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

              {/* Verification Info */}
              {selectedAchievement.verification && (
                <div className="p-4 bg-green-50 rounded-lg border border-green-200">
                  <h4 className="font-medium text-green-800 mb-2">Verification Details</h4>
                  <p className="text-sm text-green-700">
                    Verified on: {formatDate(selectedAchievement.verification.verified_at)}
                  </p>
                  {selectedAchievement.verification.comments && (
                    <p className="text-sm text-green-700 mt-1">
                      Comments: {selectedAchievement.verification.comments}
                    </p>
                  )}
                </div>
              )}

              {selectedAchievement.rejection && (
                <div className="p-4 bg-red-50 rounded-lg border border-red-200">
                  <h4 className="font-medium text-red-800 mb-2">Rejection Details</h4>
                  <p className="text-sm text-red-700">
                    Reason: {selectedAchievement.rejection.reason}
                  </p>
                </div>
              )}

              {/* Certificate Link */}
              {selectedAchievement.data?.certificate_url && (
                <div className="p-4 bg-blue-50 rounded-lg border border-blue-200">
                  <h4 className="font-medium text-blue-800 mb-2">Certificate</h4>
                  <a
                    href={getFileUrl(selectedAchievement.data.certificate_url)}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-sm text-blue-600 hover:underline"
                  >
                    View Certificate →
                  </a>
                </div>
              )}

              {/* Timestamps */}
              <div className="border-t pt-4 text-sm text-gray-500">
                <p>Created: {formatDate(selectedAchievement.created_at)}</p>
                <p>Last Updated: {formatDate(selectedAchievement.updated_at)}</p>
              </div>
            </div>
          </div>
        </div>
      )}
    </DashboardLayout>
  )
}
