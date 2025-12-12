import { useState, useEffect } from 'react'
import DashboardLayout from '../../components/DashboardLayout'
import { userService, lecturerService, studentService } from '../../services'
import { Plus, Edit, Trash2, Search, X } from 'lucide-react'

export default function Users() {
  const [users, setUsers] = useState([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const [search, setSearch] = useState('')
  const [showModal, setShowModal] = useState(false)
  const [editingUser, setEditingUser] = useState(null)
  const [saving, setSaving] = useState(false)
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    full_name: '',
    userType: 'admin',
    student_id: '',
    program_study: '',
    academic_year: new Date().getFullYear().toString(),
    lecturer_id: '',
    department: '',
  })

  useEffect(() => {
    fetchUsers()
  }, [page])

  const fetchUsers = async () => {
    setLoading(true)
    try {
      const response = await userService.getUsers(page, 10)
      if (response.status === 'success') {
        const paginationData = response.pagination?.data || response.data || {}
        setUsers(paginationData.users || [])
        setTotal(response.pagination?.total || 0)
      }
    } catch (error) {
      console.error('Error fetching users:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchLecturers = async () => {
    try {
      const response = await lecturerService.getLecturers(1, 100)
      if (response.status === 'success') {
        const paginationData = response.pagination?.data || response.data || {}
        // Lecturers data no longer needed in this component
      }
    } catch (error) {
      console.error('Error fetching lecturers:', error)
    }
  }

  const fetchStudents = async () => {
    try {
      const response = await studentService.getStudents(1, 100)
      if (response.status === 'success') {
        const paginationData = response.pagination?.data || response.data || {}
        // Students data no longer needed in this component
      }
    } catch (error) {
      console.error('Error fetching students:', error)
    }
  }

  const handleDelete = async (id) => {
    if (!id) {
      alert('Invalid user ID')
      return
    }
    if (!confirm('Are you sure you want to delete this user?')) return

    try {
      const response = await userService.deleteUser(id)
      if (response.status === 'success') {
        fetchUsers()
      } else {
        alert(response.message || 'Failed to delete user')
      }
    } catch (error) {
      console.error('Error deleting user:', error)
      alert(error.response?.data?.message || 'Failed to delete user')
    }
  }

  // Generate random ID for students/lecturers
  const generateRandomId = (prefix) => {
    const random = Math.floor(Math.random() * 900000) + 100000
    return `${prefix}${random}`
  }

  const openCreateModal = () => {
    setEditingUser(null)
    setFormData({
      username: '',
      email: '',
      password: '',
      full_name: '',
      userType: 'admin',
      student_id: generateRandomId('STU'),
      program_study: '',
      academic_year: new Date().getFullYear().toString(),
      lecturer_id: generateRandomId('LEC'),
      department: '',
      advisor_id: '',
    })
    setShowModal(true)
  }

  const openEditModal = (user) => {
    setEditingUser(user)
    const roleName = user.role?.name || ''
    let userType = 'admin'
    if (roleName === 'Mahasiswa') userType = 'student'
    else if (roleName === 'Dosen Wali') userType = 'lecturer'

    setFormData({
      username: user.username || '',
      email: user.email || '',
      password: '',
      full_name: user.full_name || '',
      userType: userType,
      student_id: generateRandomId('STU'),
      program_study: '',
      academic_year: new Date().getFullYear().toString(),
      lecturer_id: generateRandomId('LEC'),
      department: '',
    })
    setShowModal(true)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setSaving(true)

    try {
      let userData = {
        username: formData.username,
        email: formData.email,
        full_name: formData.full_name,
      }

      if (!editingUser) {
        userData.password = formData.password
      } else if (formData.password) {
        userData.password = formData.password
      }

      // Add role-specific fields
      if (formData.userType === 'student') {
        userData.student_id = formData.student_id
        userData.program_study = formData.program_study
        userData.academic_year = formData.academic_year
      } else if (formData.userType === 'lecturer') {
        userData.lecturer_id = formData.lecturer_id
        userData.department = formData.department
      } else if (formData.userType === 'admin') {
        userData.role_name = 'Admin'
      }

      if (editingUser) {
        await userService.updateUser(editingUser.id, userData)
      } else {
        await userService.createUser(userData)
      }

      setShowModal(false)
      fetchUsers()
    } catch (error) {
      console.error('Error saving user:', error)
      alert(error.response?.data?.message || 'Failed to save user')
    } finally {
      setSaving(false)
    }
  }

  const filteredUsers = users.filter(user =>
    user.full_name?.toLowerCase().includes(search.toLowerCase()) ||
    user.username?.toLowerCase().includes(search.toLowerCase()) ||
    user.email?.toLowerCase().includes(search.toLowerCase())
  )

  const totalPages = Math.ceil(total / 10)

  return (
    <DashboardLayout title="User Management">
      <div className="space-y-6">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div className="relative flex-1 max-w-md">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
            <input
              type="text"
              placeholder="Search users..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="input pl-10 w-full"
            />
          </div>
          <button onClick={openCreateModal} className="btn btn-primary whitespace-nowrap">
            <Plus className="w-5 h-5 mr-2" />
            Add User
          </button>
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
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
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
                    <td colSpan="5" className="px-6 py-4 text-center text-gray-500">No users found</td>
                  </tr>
                ) : filteredUsers.map((user) => (
                  <tr key={user.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        <div className="h-10 w-10 flex-shrink-0">
                          <div className="h-10 w-10 rounded-full bg-primary-100 flex items-center justify-center">
                            <span className="text-primary-600 font-medium">
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
                    <td className="px-6 py-4 whitespace-nowrap">
                      {user.is_active ? (
                        <span className="badge badge-success">Active</span>
                      ) : (
                        <span className="badge badge-danger">Inactive</span>
                      )}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <button
                        onClick={() => openEditModal(user)}
                        className="text-primary-600 hover:text-primary-900 mr-3"
                        title="Edit User"
                      >
                        <Edit className="w-5 h-5" />
                      </button>
                      <button
                        onClick={() => handleDelete(user.id)}
                        className="text-red-600 hover:text-red-900"
                        title="Delete User"
                      >
                        <Trash2 className="w-5 h-5" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
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
        </div>
      </div>

      {/* Create/Edit User Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-xl max-w-md w-full max-h-[90vh] overflow-y-auto">
            <div className="flex items-center justify-between p-6 border-b">
              <h2 className="text-xl font-bold text-gray-900">
                {editingUser ? 'Edit User' : 'Create New User'}
              </h2>
              <button onClick={() => setShowModal(false)} className="text-gray-400 hover:text-gray-600">
                <X className="w-6 h-6" />
              </button>
            </div>

            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              {/* User Type */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">User Type *</label>
                <div className="flex gap-4">
                  {['admin', 'student', 'lecturer'].map((type) => (
                    <label key={type} className="flex items-center">
                      <input
                        type="radio"
                        name="userType"
                        value={type}
                        checked={formData.userType === type}
                        onChange={(e) => setFormData({ ...formData, userType: e.target.value })}
                        className="mr-2"
                        disabled={!!editingUser}
                      />
                      <span className="capitalize">{type === 'student' ? 'Mahasiswa' : type === 'lecturer' ? 'Dosen Wali' : type}</span>
                    </label>
                  ))}
                </div>
              </div>

              {/* Basic Info */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Full Name *</label>
                  <input
                    type="text"
                    value={formData.full_name}
                    onChange={(e) => setFormData({ ...formData, full_name: e.target.value })}
                    className="input w-full"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Username *</label>
                  <input
                    type="text"
                    value={formData.username}
                    onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                    className="input w-full"
                    required
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Email *</label>
                <input
                  type="email"
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  className="input w-full"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Password {editingUser ? '(leave empty to keep current)' : '*'}
                </label>
                <input
                  type="password"
                  value={formData.password}
                  onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                  className="input w-full"
                  required={!editingUser}
                />
              </div>

              {/* Student-specific fields */}
              {formData.userType === 'student' && (
                <div className="border-t pt-4 mt-4">
                  <h3 className="font-medium text-gray-900 mb-3">Student Information</h3>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Student ID *</label>
                      <div className="flex gap-2">
                        <input
                          type="text"
                          value={formData.student_id}
                          onChange={(e) => setFormData({ ...formData, student_id: e.target.value })}
                          className="input w-full"
                          required
                        />
                        <button
                          type="button"
                          onClick={() => setFormData({ ...formData, student_id: generateRandomId('STU') })}
                          className="btn btn-secondary btn-sm whitespace-nowrap"
                          title="Generate Random ID"
                        >
                          🔄
                        </button>
                      </div>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Academic Year *</label>
                      <input
                        type="text"
                        value={formData.academic_year}
                        onChange={(e) => setFormData({ ...formData, academic_year: e.target.value })}
                        className="input w-full"
                        required
                      />
                    </div>
                  </div>
                  <div className="mt-3">
                    <label className="block text-sm font-medium text-gray-700 mb-1">Program Study *</label>
                    <input
                      type="text"
                      value={formData.program_study}
                      onChange={(e) => setFormData({ ...formData, program_study: e.target.value })}
                      className="input w-full"
                      required
                    />
                  </div>
                </div>
              )}

              {/* Lecturer-specific fields */}
              {formData.userType === 'lecturer' && (
                <div className="border-t pt-4 mt-4">
                  <h3 className="font-medium text-gray-900 mb-3">Lecturer Information</h3>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Lecturer ID *</label>
                      <div className="flex gap-2">
                        <input
                          type="text"
                          value={formData.lecturer_id}
                          onChange={(e) => setFormData({ ...formData, lecturer_id: e.target.value })}
                          className="input w-full"
                          required
                        />
                        <button
                          type="button"
                          onClick={() => setFormData({ ...formData, lecturer_id: generateRandomId('LEC') })}
                          className="btn btn-secondary btn-sm whitespace-nowrap"
                          title="Generate Random ID"
                        >
                          🔄
                        </button>
                      </div>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Department *</label>
                      <input
                        type="text"
                        value={formData.department}
                        onChange={(e) => setFormData({ ...formData, department: e.target.value })}
                        className="input w-full"
                        required
                      />
                    </div>
                  </div>
                </div>
              )}

              <div className="flex gap-3 pt-4">
                <button type="button" onClick={() => setShowModal(false)} className="btn btn-secondary flex-1">
                  Cancel
                </button>
                <button type="submit" className="btn btn-primary flex-1" disabled={saving}>
                  {saving ? 'Saving...' : (editingUser ? 'Update User' : 'Create User')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </DashboardLayout>
  )
}
