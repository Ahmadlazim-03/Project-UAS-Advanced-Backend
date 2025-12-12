import { useState, useEffect } from 'react'
import DashboardLayout from '../../components/DashboardLayout'
import { userService } from '../../services'
import { Search, RotateCcw, Trash2, AlertTriangle, CheckCircle, X } from 'lucide-react'

export default function DeletedUsers() {
  const [deletedUsers, setDeletedUsers] = useState([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const [search, setSearch] = useState('')
  const [showConfirmModal, setShowConfirmModal] = useState(false)
  const [selectedUser, setSelectedUser] = useState(null)
  const [actionType, setActionType] = useState('') // 'restore' or 'hardDelete'
  const [processing, setProcessing] = useState(false)
  const [showSuccessToast, setShowSuccessToast] = useState(false)
  const [successMessage, setSuccessMessage] = useState('')

  useEffect(() => {
    fetchDeletedUsers()
  }, [page])

  const fetchDeletedUsers = async () => {
    setLoading(true)
    try {
      const response = await userService.getDeletedUsers(page, 10)
      if (response.status === 'success') {
        const paginationData = response.pagination?.data || response.data || {}
        setDeletedUsers(paginationData.users || [])
        setTotal(response.pagination?.total || 0)
      }
    } catch (error) {
      console.error('Error fetching deleted users:', error)
    } finally {
      setLoading(false)
    }
  }

  const openConfirmModal = (user, action) => {
    setSelectedUser(user)
    setActionType(action)
    setShowConfirmModal(true)
  }

  const handleRestore = async () => {
    if (!selectedUser) return
    
    setProcessing(true)
    try {
      const response = await userService.restoreUser(selectedUser.id)
      if (response.status === 'success') {
        setShowConfirmModal(false)
        fetchDeletedUsers()
        setSuccessMessage(`User "${selectedUser.full_name}" has been restored successfully!`)
        setShowSuccessToast(true)
        setTimeout(() => setShowSuccessToast(false), 3000)
      } else {
        alert(response.message || 'Failed to restore user')
      }
    } catch (error) {
      console.error('Error restoring user:', error)
      alert(error.response?.data?.message || 'Failed to restore user')
    } finally {
      setProcessing(false)
    }
  }

  const handleHardDelete = async () => {
    if (!selectedUser) return
    
    setProcessing(true)
    try {
      const response = await userService.hardDeleteUser(selectedUser.id)
      if (response.status === 'success') {
        setShowConfirmModal(false)
        fetchDeletedUsers()
        setSuccessMessage(`User "${selectedUser.full_name}" has been permanently deleted!`)
        setShowSuccessToast(true)
        setTimeout(() => setShowSuccessToast(false), 3000)
      } else {
        alert(response.message || 'Failed to delete user permanently')
      }
    } catch (error) {
      console.error('Error deleting user:', error)
      alert(error.response?.data?.message || 'Failed to delete user permanently')
    } finally {
      setProcessing(false)
    }
  }

  const handleConfirm = () => {
    if (actionType === 'restore') {
      handleRestore()
    } else if (actionType === 'hardDelete') {
      handleHardDelete()
    }
  }

  const filteredUsers = deletedUsers.filter(user =>
    user.full_name?.toLowerCase().includes(search.toLowerCase()) ||
    user.username?.toLowerCase().includes(search.toLowerCase()) ||
    user.email?.toLowerCase().includes(search.toLowerCase())
  )

  const totalPages = Math.ceil(total / 10)

  return (
    <DashboardLayout title="Deleted Users">
      {/* Success Toast */}
      {showSuccessToast && (
        <div className="fixed top-4 right-4 z-50 animate-fade-in-down">
          <div className="bg-white rounded-lg shadow-lg border-l-4 border-green-500 p-4 flex items-start max-w-md">
            <CheckCircle className="h-6 w-6 text-green-500 flex-shrink-0" />
            <div className="ml-3 flex-1">
              <p className="text-sm font-medium text-gray-900">Success!</p>
              <p className="text-sm text-gray-600 mt-1">{successMessage}</p>
            </div>
            <button
              onClick={() => setShowSuccessToast(false)}
              className="ml-4 flex-shrink-0 text-gray-400 hover:text-gray-500"
            >
              <X className="h-5 w-5" />
            </button>
          </div>
        </div>
      )}

      <div className="space-y-6">
        {/* Info Banner */}
        <div className="bg-yellow-50 border-l-4 border-yellow-400 p-4 rounded-lg">
          <div className="flex items-start">
            <AlertTriangle className="h-5 w-5 text-yellow-400 mt-0.5" />
            <div className="ml-3">
              <h3 className="text-sm font-medium text-yellow-800">Soft-Deleted Users</h3>
              <p className="text-sm text-yellow-700 mt-1">
                These users have been soft-deleted. You can restore them or permanently delete them from the database.
                Soft-deleted users will not appear in the Advisor Management list.
              </p>
            </div>
          </div>
        </div>

        {/* Search */}
        <div className="card">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
            <input
              type="text"
              placeholder="Search deleted users..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="input pl-10 w-full"
            />
          </div>
        </div>

        {/* Users Table */}
        <div className="card overflow-hidden">
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">User</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Role</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Deleted At</th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {loading ? (
                  <tr>
                    <td colSpan="5" className="px-6 py-4 text-center">
                      <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
                    </td>
                  </tr>
                ) : filteredUsers.length === 0 ? (
                  <tr>
                    <td colSpan="5" className="px-6 py-4 text-center text-gray-500">
                      {search ? 'No deleted users found matching your search' : 'No deleted users'}
                    </td>
                  </tr>
                ) : filteredUsers.map((user) => (
                  <tr key={user.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        <div className="h-10 w-10 flex-shrink-0">
                          <div className="h-10 w-10 rounded-full bg-gray-200 flex items-center justify-center">
                            <span className="text-gray-600 font-medium">
                              {user.full_name?.charAt(0)?.toUpperCase() || '?'}
                            </span>
                          </div>
                        </div>
                        <div className="ml-4">
                          <div className="text-sm font-medium text-gray-900">{user.full_name}</div>
                          <div className="text-sm text-gray-500">@{user.username}</div>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {user.email}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`badge ${user.role?.name === 'Admin' ? 'badge-danger' :
                        user.role?.name === 'Dosen Wali' ? 'badge-info' :
                          'badge-success'
                        }`}>{user.role?.name || 'No Role'}</span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {user.deleted_at ? new Date(user.deleted_at).toLocaleDateString() : '-'}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-2">
                      <button
                        onClick={() => openConfirmModal(user, 'restore')}
                        className="text-green-600 hover:text-green-900 inline-flex items-center"
                        title="Restore User"
                      >
                        <RotateCcw className="w-5 h-5 mr-1" />
                        Restore
                      </button>
                      <button
                        onClick={() => openConfirmModal(user, 'hardDelete')}
                        className="text-red-600 hover:text-red-900 inline-flex items-center"
                        title="Permanently Delete"
                      >
                        <Trash2 className="w-5 h-5 mr-1" />
                        Delete Forever
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          {total > 10 && (
            <div className="px-6 py-4 border-t border-gray-200 flex items-center justify-between">
              <div className="text-sm text-gray-700">
                Showing {((page - 1) * 10) + 1} to {Math.min(page * 10, total)} of {total} users
              </div>
              <div className="flex gap-2">
                <button
                  onClick={() => setPage(p => Math.max(1, p - 1))}
                  disabled={page === 1}
                  className="btn btn-secondary btn-sm"
                >
                  Previous
                </button>
                <button
                  onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                  disabled={page >= totalPages}
                  className="btn btn-secondary btn-sm"
                >
                  Next
                </button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Confirmation Modal */}
      {showConfirmModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-xl max-w-md w-full">
            <div className="flex items-center justify-between p-6 border-b">
              <h2 className="text-xl font-bold text-gray-900">
                {actionType === 'restore' ? 'Restore User' : 'Permanently Delete User'}
              </h2>
              <button
                onClick={() => setShowConfirmModal(false)}
                className="text-gray-400 hover:text-gray-600"
                disabled={processing}
              >
                <X className="w-6 h-6" />
              </button>
            </div>

            <div className="p-6">
              {actionType === 'restore' ? (
                <div className="space-y-4">
                  <div className="flex items-start">
                    <RotateCcw className="h-6 w-6 text-green-500 mt-0.5" />
                    <div className="ml-3">
                      <p className="text-sm text-gray-700">
                        Are you sure you want to restore <strong>{selectedUser?.full_name}</strong>?
                      </p>
                      <p className="text-sm text-gray-500 mt-2">
                        This user will be reactivated and will appear in all user lists including Advisor Management.
                      </p>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="space-y-4">
                  <div className="flex items-start">
                    <AlertTriangle className="h-6 w-6 text-red-500 mt-0.5" />
                    <div className="ml-3">
                      <p className="text-sm text-gray-700">
                        Are you sure you want to <strong className="text-red-600">permanently delete</strong> <strong>{selectedUser?.full_name}</strong>?
                      </p>
                      <p className="text-sm text-red-600 mt-2 font-medium">
                        This action cannot be undone! The user and all associated data will be permanently removed from the database.
                      </p>
                    </div>
                  </div>
                </div>
              )}
            </div>

            <div className="flex gap-3 p-6 border-t">
              <button
                onClick={() => setShowConfirmModal(false)}
                className="btn btn-secondary flex-1"
                disabled={processing}
              >
                Cancel
              </button>
              <button
                onClick={handleConfirm}
                className={`btn flex-1 ${actionType === 'restore' ? 'btn-primary' : 'btn-danger'}`}
                disabled={processing}
              >
                {processing ? 'Processing...' : actionType === 'restore' ? 'Restore User' : 'Delete Forever'}
              </button>
            </div>
          </div>
        </div>
      )}
    </DashboardLayout>
  )
}
