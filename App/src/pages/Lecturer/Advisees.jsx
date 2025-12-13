import { useState, useEffect } from 'react'
import DashboardLayout from '../../components/DashboardLayout'
import { lecturerService } from '../../services'
import {
  Users,
  Search,
  Mail,
  Phone,
  Calendar,
  Award,
  FileText,
  X,
  GraduationCap,
  CheckCircle,
  Clock,
  XCircle,
} from 'lucide-react'

export default function Advisees() {
  const [students, setStudents] = useState([])
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedStudent, setSelectedStudent] = useState(null)
  const [studentReport, setStudentReport] = useState(null)
  const [loadingReport, setLoadingReport] = useState(false)

  useEffect(() => {
    fetchAdvisees()
  }, [])

  const fetchAdvisees = async () => {
    try {
      setLoading(true)
      const response = await lecturerService.getAdvisees()
      console.log('Advisees response:', response)
      if (response.status === 'success') {
        const studentsList = response.data || []
        console.log('Students list:', studentsList)
        setStudents(studentsList)
      }
    } catch (error) {
      console.error('Error fetching advisees:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchStudentReport = async (studentUserId) => {
    setLoadingReport(true)
    try {
      console.log('Fetching report for student user_id:', studentUserId)
      const response = await lecturerService.getStudentReport(studentUserId)
      console.log('Student report response:', response)
      if (response.status === 'success') {
        setStudentReport(response.data)
      } else {
        console.error('Failed to fetch student report:', response)
        alert(`Failed to load student report: ${response.message || 'Unknown error'}`)
      }
    } catch (error) {
      console.error('Error fetching student report:', error)
      console.error('Error details:', error.response?.data)
      alert(`Failed to load student report: ${error.response?.data?.message || error.message}`)
    } finally {
      setLoadingReport(false)
    }
  }

  const handleStudentClick = (student) => {
    setSelectedStudent(student)
    fetchStudentReport(student.user_id)
  }

  const closeModal = () => {
    setSelectedStudent(null)
    setStudentReport(null)
  }

  const filteredStudents = students.filter(student =>
    student.user?.full_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
    student.student_id?.toLowerCase().includes(searchQuery.toLowerCase()) ||
    student.program_study?.toLowerCase().includes(searchQuery.toLowerCase())
  )

  if (loading) {
    return (
      <DashboardLayout title="My Advisees">
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
        </div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout title="My Advisees">
      <div className="space-y-6">
        {/* Header */}
        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h2 className="text-2xl font-bold text-gray-900">My Advisees</h2>
              <p className="text-gray-600 mt-1">
                Total: {students.length} student{students.length !== 1 ? 's' : ''}
              </p>
            </div>
            <div className="bg-primary-500 p-3 rounded-lg">
              <Users className="w-8 h-8 text-white" />
            </div>
          </div>

          {/* Search */}
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
            <input
              type="text"
              placeholder="Search by name, NIM, or program..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="input pl-10"
            />
          </div>
        </div>

        {/* Students List */}
        {filteredStudents.length === 0 ? (
          <div className="card text-center py-12">
            <Users className="w-16 h-16 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-600">
              {searchQuery ? 'No students found' : 'No advisees assigned yet'}
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredStudents.map((student) => (
              <div
                key={student.id}
                onClick={() => handleStudentClick(student)}
                className="card hover:shadow-xl transition-all cursor-pointer border-2 border-transparent hover:border-primary-500"
              >
                <div className="flex items-start justify-between mb-4">
                  <div className="flex-1">
                    <h3 className="font-semibold text-lg text-gray-900 mb-1">
                      {student.user?.full_name || 'Unknown'}
                    </h3>
                    <p className="text-sm text-gray-500">{student.student_id}</p>
                  </div>
                  <div className="bg-primary-100 p-2 rounded-lg">
                    <GraduationCap className="w-5 h-5 text-primary-600" />
                  </div>
                </div>

                <div className="space-y-2 text-sm">
                  <div className="flex items-center text-gray-600">
                    <Mail className="w-4 h-4 mr-2" />
                    <span className="truncate">{student.user?.email || 'N/A'}</span>
                  </div>
                  <div className="flex items-center text-gray-600">
                    <FileText className="w-4 h-4 mr-2" />
                    <span>{student.program_study || 'N/A'}</span>
                  </div>
                  <div className="flex items-center text-gray-600">
                    <Calendar className="w-4 h-4 mr-2" />
                    <span>Angkatan {student.academic_year || 'N/A'}</span>
                  </div>
                </div>

                <div className="mt-4 pt-4 border-t border-gray-200">
                  <button className="text-primary-600 hover:text-primary-700 font-medium text-sm flex items-center">
                    <Award className="w-4 h-4 mr-1" />
                    View Details
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Student Detail Modal */}
      {selectedStudent && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-xl max-w-3xl w-full max-h-[90vh] overflow-y-auto">
            {/* Modal Header */}
            <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
              <div>
                <h2 className="text-2xl font-bold text-gray-900">
                  {selectedStudent.user?.full_name}
                </h2>
                <p className="text-gray-600">{selectedStudent.student_id}</p>
              </div>
              <button
                onClick={closeModal}
                className="text-gray-500 hover:text-gray-700 p-2 hover:bg-gray-100 rounded-lg transition-colors"
              >
                <X className="w-6 h-6" />
              </button>
            </div>

            {/* Modal Content */}
            <div className="p-6">
              {loadingReport ? (
                <div className="flex items-center justify-center py-12">
                  <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
                </div>
              ) : studentReport ? (
                <div className="space-y-6">
                  {/* Student Info */}
                  <div className="card bg-gray-50">
                    <h3 className="font-semibold text-lg text-gray-900 mb-4">Personal Information</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                      <div>
                        <p className="text-gray-600">Email</p>
                        <p className="font-medium text-gray-900">{selectedStudent.user?.email}</p>
                      </div>
                      <div>
                        <p className="text-gray-600">Program Study</p>
                        <p className="font-medium text-gray-900">{selectedStudent.program_study}</p>
                      </div>
                      <div>
                        <p className="text-gray-600">Academic Year</p>
                        <p className="font-medium text-gray-900">{selectedStudent.academic_year}</p>
                      </div>
                      <div>
                        <p className="text-gray-600">Student ID</p>
                        <p className="font-medium text-gray-900">{selectedStudent.student_id}</p>
                      </div>
                    </div>
                  </div>

                  {/* Achievement Summary */}
                  <div>
                    <h3 className="font-semibold text-lg text-gray-900 mb-4">Achievement Summary</h3>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                      <div className="card bg-blue-50 border-blue-200">
                        <div className="flex items-center justify-between mb-2">
                          <Award className="w-5 h-5 text-blue-600" />
                        </div>
                        <p className="text-2xl font-bold text-blue-600">
                          {studentReport.summary?.total_achievements || 0}
                        </p>
                        <p className="text-xs text-gray-600 mt-1">Total Achievements</p>
                      </div>

                      <div className="card bg-green-50 border-green-200">
                        <div className="flex items-center justify-between mb-2">
                          <CheckCircle className="w-5 h-5 text-green-600" />
                        </div>
                        <p className="text-2xl font-bold text-green-600">
                          {studentReport.summary?.verified_achievements || 0}
                        </p>
                        <p className="text-xs text-gray-600 mt-1">Verified</p>
                      </div>

                      <div className="card bg-yellow-50 border-yellow-200">
                        <div className="flex items-center justify-between mb-2">
                          <Clock className="w-5 h-5 text-yellow-600" />
                        </div>
                        <p className="text-2xl font-bold text-yellow-600">
                          {studentReport.summary?.pending_achievements || 0}
                        </p>
                        <p className="text-xs text-gray-600 mt-1">Pending</p>
                      </div>

                      <div className="card bg-red-50 border-red-200">
                        <div className="flex items-center justify-between mb-2">
                          <XCircle className="w-5 h-5 text-red-600" />
                        </div>
                        <p className="text-2xl font-bold text-red-600">
                          {studentReport.summary?.rejected_achievements || 0}
                        </p>
                        <p className="text-xs text-gray-600 mt-1">Rejected</p>
                      </div>
                    </div>
                  </div>

                  {/* Achievement Types */}
                  {studentReport.achievements_by_type && Object.keys(studentReport.achievements_by_type).length > 0 && (
                    <div>
                      <h3 className="font-semibold text-lg text-gray-900 mb-4">Achievements by Type</h3>
                      <div className="card bg-gray-50">
                        <div className="space-y-3">
                          {Object.entries(studentReport.achievements_by_type).map(([type, count]) => (
                            <div key={type} className="flex items-center justify-between">
                              <span className="text-gray-700 capitalize">{type}</span>
                              <span className="font-semibold text-primary-600">{count}</span>
                            </div>
                          ))}
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              ) : (
                <div className="text-center py-12">
                  <p className="text-gray-600">No data available</p>
                </div>
              )}
            </div>

            {/* Modal Footer */}
            <div className="sticky bottom-0 bg-gray-50 border-t border-gray-200 px-6 py-4 flex justify-end">
              <button
                onClick={closeModal}
                className="btn-secondary"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}
    </DashboardLayout>
  )
}
