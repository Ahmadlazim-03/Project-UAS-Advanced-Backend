import { useState, useEffect } from 'react'
import DashboardLayout from '../../components/DashboardLayout'
import { studentService, lecturerService } from '../../services'
import { Search, X, UserPlus, Users, GraduationCap, CheckCircle } from 'lucide-react'

export default function Advisors() {
    const [students, setStudents] = useState([])
    const [lecturers, setLecturers] = useState([])
    const [loading, setLoading] = useState(true)
    const [showAdvisorModal, setShowAdvisorModal] = useState(false)
    const [selectedStudent, setSelectedStudent] = useState(null)
    const [saving, setSaving] = useState(false)
    const [advisorId, setAdvisorId] = useState('')
    const [searchTerm, setSearchTerm] = useState('')
    const [filterStatus, setFilterStatus] = useState('all') // all, assigned, unassigned
    const [showSuccessToast, setShowSuccessToast] = useState(false)
    const [successMessage, setSuccessMessage] = useState('')

    useEffect(() => {
        fetchData()
    }, [])

    const fetchData = async () => {
        setLoading(true)
        try {
            const [studentsRes, lecturersRes] = await Promise.all([
                studentService.getStudents(1, 100),
                lecturerService.getLecturers(1, 100)
            ])

            if (studentsRes.status === 'success') {
                const paginationData = studentsRes.pagination?.data || studentsRes.data || {}
                setStudents(paginationData.students || [])
            }

            if (lecturersRes.status === 'success') {
                const paginationData = lecturersRes.pagination?.data || lecturersRes.data || {}
                setLecturers(paginationData.lecturers || [])
            }
        } catch (error) {
            console.error('Error fetching data:', error)
        } finally {
            setLoading(false)
        }
    }

    const openAdvisorModal = (student) => {
        setSelectedStudent(student)
        setAdvisorId(student.advisor_id || '')
        setShowAdvisorModal(true)
    }

    const handleAssignAdvisor = async () => {
        if (!selectedStudent || !advisorId) {
            alert('Please select an advisor')
            return
        }

        setSaving(true)
        try {
            // Use user_id instead of id for the API call
            const studentUserId = selectedStudent.user_id
            if (!studentUserId) {
                throw new Error('Student user_id not found')
            }

            const response = await studentService.assignAdvisor(studentUserId, advisorId)
            if (response.status === 'success') {
                setShowAdvisorModal(false)
                fetchData()
                setSuccessMessage(`Advisor successfully assigned to ${selectedStudent?.user?.full_name}!`)
                setShowSuccessToast(true)
                setTimeout(() => setShowSuccessToast(false), 3000)
            } else {
                alert(response.message || 'Failed to assign advisor')
            }
        } catch (error) {
            console.error('Error assigning advisor:', error)
            const errorMessage = error.response?.data?.message || error.message || 'Failed to assign advisor'
            alert(errorMessage)
        } finally {
            setSaving(false)
        }
    }

    // Filter students based on search and status
    const filteredStudents = students.filter(student => {
        const matchesSearch =
            student.user?.full_name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
            student.student_id?.toLowerCase().includes(searchTerm.toLowerCase()) ||
            student.program_study?.toLowerCase().includes(searchTerm.toLowerCase())

        const matchesFilter =
            filterStatus === 'all' ||
            (filterStatus === 'assigned' && student.advisor) ||
            (filterStatus === 'unassigned' && !student.advisor)

        return matchesSearch && matchesFilter
    })

    // Statistics
    const totalStudents = students.length
    const assignedCount = students.filter(s => s.advisor).length
    const unassignedCount = students.filter(s => !s.advisor).length

    return (
        <DashboardLayout title="Advisor Management">
            {/* Success Toast Notification */}
            {showSuccessToast && (
                <div className="fixed top-4 right-4 z-50 animate-fade-in-down">
                    <div className="bg-white rounded-lg shadow-lg border-l-4 border-green-500 p-4 flex items-start max-w-md">
                        <div className="flex-shrink-0">
                            <CheckCircle className="h-6 w-6 text-green-500" />
                        </div>
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
                {/* Statistics Cards */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                    <div className="card bg-gradient-to-r from-blue-500 to-blue-600 text-white">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-blue-100 text-sm">Total Students</p>
                                <p className="text-3xl font-bold">{totalStudents}</p>
                            </div>
                            <GraduationCap className="w-12 h-12 opacity-80" />
                        </div>
                    </div>

                    <div className="card bg-gradient-to-r from-green-500 to-green-600 text-white">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-green-100 text-sm">With Advisor</p>
                                <p className="text-3xl font-bold">{assignedCount}</p>
                            </div>
                            <UserPlus className="w-12 h-12 opacity-80" />
                        </div>
                    </div>

                    <div className="card bg-gradient-to-r from-amber-500 to-amber-600 text-white">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-amber-100 text-sm">Without Advisor</p>
                                <p className="text-3xl font-bold">{unassignedCount}</p>
                            </div>
                            <Users className="w-12 h-12 opacity-80" />
                        </div>
                    </div>
                </div>

                {/* Filters */}
                <div className="card">
                    <div className="flex flex-col md:flex-row md:items-center gap-4">
                        <div className="relative flex-1">
                            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
                            <input
                                type="text"
                                placeholder="Search by name, student ID, or program..."
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                                className="input pl-10 w-full"
                            />
                        </div>
                        <div className="flex gap-2">
                            <button
                                onClick={() => setFilterStatus('all')}
                                className={`btn btn-sm ${filterStatus === 'all' ? 'btn-primary' : 'btn-secondary'}`}
                            >
                                All ({totalStudents})
                            </button>
                            <button
                                onClick={() => setFilterStatus('assigned')}
                                className={`btn btn-sm ${filterStatus === 'assigned' ? 'btn-primary' : 'btn-secondary'}`}
                            >
                                Assigned ({assignedCount})
                            </button>
                            <button
                                onClick={() => setFilterStatus('unassigned')}
                                className={`btn btn-sm ${filterStatus === 'unassigned' ? 'btn-primary' : 'btn-secondary'}`}
                            >
                                Unassigned ({unassignedCount})
                            </button>
                        </div>
                    </div>
                </div>

                {/* Students Table */}
                <div className="card overflow-hidden">
                    <div className="overflow-x-auto">
                        <table className="min-w-full divide-y divide-gray-200">
                            <thead className="bg-gray-50">
                                <tr>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Student</th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Student ID</th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Program</th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Academic Year</th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Current Advisor</th>
                                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Action</th>
                                </tr>
                            </thead>
                            <tbody className="bg-white divide-y divide-gray-200">
                                {loading ? (
                                    <tr>
                                        <td colSpan="6" className="px-6 py-8 text-center">
                                            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
                                        </td>
                                    </tr>
                                ) : filteredStudents.length === 0 ? (
                                    <tr>
                                        <td colSpan="6" className="px-6 py-8 text-center text-gray-500">
                                            No students found
                                        </td>
                                    </tr>
                                ) : filteredStudents.map((student) => (
                                    <tr key={student.id} className="hover:bg-gray-50">
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <div className="flex items-center">
                                                <div className="h-10 w-10 flex-shrink-0">
                                                    <div className="h-10 w-10 rounded-full bg-primary-100 flex items-center justify-center">
                                                        <span className="text-primary-600 font-medium">
                                                            {student.user?.full_name?.charAt(0)?.toUpperCase() || '?'}
                                                        </span>
                                                    </div>
                                                </div>
                                                <div className="ml-4">
                                                    <div className="text-sm font-medium text-gray-900">
                                                        {student.user?.full_name || 'Unknown'}
                                                    </div>
                                                    <div className="text-sm text-gray-500">
                                                        {student.user?.email}
                                                    </div>
                                                </div>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <span className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">
                                                {student.student_id}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                                            {student.program_study || '-'}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                                            {student.academic_year || '-'}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            {student.advisor ? (
                                                <div className="flex items-center">
                                                    <span className="badge badge-success mr-2">
                                                        {student.advisor.user?.full_name || 'Assigned'}
                                                    </span>
                                                    <span className="text-xs text-gray-500">
                                                        {student.advisor.department}
                                                    </span>
                                                </div>
                                            ) : (
                                                <span className="badge badge-warning">No Advisor</span>
                                            )}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-right">
                                            <button
                                                onClick={() => openAdvisorModal(student)}
                                                className={`btn btn-sm ${student.advisor ? 'btn-secondary' : 'btn-primary'}`}
                                            >
                                                <UserPlus className="w-4 h-4 mr-1" />
                                                {student.advisor ? 'Change' : 'Assign'}
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>

                {/* Available Lecturers Info */}
                <div className="card">
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Available Advisors (Dosen Wali)</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                        {lecturers.map((lecturer) => {
                            const adviseeCount = students.filter(s => s.advisor_id === lecturer.id).length
                            return (
                                <div key={lecturer.id} className="p-4 bg-gray-50 rounded-lg border border-gray-200">
                                    <div className="flex items-center space-x-3">
                                        <div className="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center">
                                            <span className="text-blue-600 font-medium">
                                                {lecturer.user?.full_name?.charAt(0)?.toUpperCase() || '?'}
                                            </span>
                                        </div>
                                        <div className="flex-1 min-w-0">
                                            <p className="text-sm font-medium text-gray-900 truncate">
                                                {lecturer.user?.full_name || lecturer.lecturer_id}
                                            </p>
                                            <p className="text-xs text-gray-500">{lecturer.department}</p>
                                        </div>
                                        <div className="text-right">
                                            <span className="text-lg font-bold text-primary-600">{adviseeCount}</span>
                                            <p className="text-xs text-gray-500">students</p>
                                        </div>
                                    </div>
                                </div>
                            )
                        })}
                    </div>
                </div>
            </div>

            {/* Assign Advisor Modal */}
            {showAdvisorModal && (
                <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
                    <div className="bg-white rounded-lg shadow-xl max-w-md w-full">
                        <div className="flex items-center justify-between p-6 border-b">
                            <h2 className="text-xl font-bold text-gray-900">
                                {selectedStudent?.advisor ? 'Change Advisor' : 'Assign Advisor'}
                            </h2>
                            <button onClick={() => setShowAdvisorModal(false)} className="text-gray-400 hover:text-gray-600">
                                <X className="w-6 h-6" />
                            </button>
                        </div>

                        <div className="p-6 space-y-4">
                            <div className="bg-blue-50 p-4 rounded-lg">
                                <p className="text-sm text-blue-800">
                                    <strong>Student:</strong> {selectedStudent?.user?.full_name || 'Unknown'}
                                </p>
                                <p className="text-sm text-blue-600">
                                    ID: {selectedStudent?.student_id}
                                </p>
                                <p className="text-sm text-blue-600">
                                    Program: {selectedStudent?.program_study}
                                </p>
                            </div>

                            {selectedStudent?.advisor && (
                                <div className="bg-green-50 p-4 rounded-lg">
                                    <p className="text-sm text-green-800">
                                        <strong>Current Advisor:</strong> {selectedStudent.advisor.user?.full_name}
                                    </p>
                                    <p className="text-sm text-green-600">
                                        Department: {selectedStudent.advisor.department}
                                    </p>
                                </div>
                            )}

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-2">
                                    Select New Advisor (Dosen Wali) *
                                </label>
                                <select
                                    value={advisorId}
                                    onChange={(e) => setAdvisorId(e.target.value)}
                                    className="input w-full"
                                    required
                                >
                                    <option value="">-- Select Advisor --</option>
                                    {lecturers.map((lecturer) => (
                                        <option key={lecturer.id} value={lecturer.id}>
                                            {lecturer.user?.full_name || lecturer.lecturer_id} - {lecturer.department}
                                        </option>
                                    ))}
                                </select>
                            </div>

                            <div className="flex gap-3 pt-4">
                                <button onClick={() => setShowAdvisorModal(false)} className="btn btn-secondary flex-1">
                                    Cancel
                                </button>
                                <button
                                    onClick={handleAssignAdvisor}
                                    className="btn btn-primary flex-1"
                                    disabled={saving || !advisorId}
                                >
                                    {saving ? 'Saving...' : (selectedStudent?.advisor ? 'Change Advisor' : 'Assign Advisor')}
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </DashboardLayout>
    )
}
