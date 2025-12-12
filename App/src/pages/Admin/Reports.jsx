import { useState, useEffect } from 'react'
import DashboardLayout from '../../components/DashboardLayout'
import { reportService, studentService } from '../../services'
import {
  Award,
  Users,
  GraduationCap,
  TrendingUp,
  Search,
  FileText,
  BarChart3,
  CheckCircle,
  Clock,
  XCircle,
  Trophy
} from 'lucide-react'

export default function Reports() {
  const [stats, setStats] = useState(null)
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [students, setStudents] = useState([])
  const [selectedStudent, setSelectedStudent] = useState(null)
  const [studentReport, setStudentReport] = useState(null)
  const [loadingReport, setLoadingReport] = useState(false)

  useEffect(() => {
    fetchStatistics()
    fetchStudents()
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

  const fetchStudents = async () => {
    try {
      const response = await studentService.getStudents(1, 100)
      if (response.status === 'success') {
        // Backend returns: { status, pagination: { data: { students: [...] } } }
        const paginationData = response.pagination?.data || response.data || {}
        setStudents(paginationData.students || [])
      }
    } catch (error) {
      console.error('Error fetching students:', error)
    }
  }

  const fetchStudentReport = async (studentId) => {
    setLoadingReport(true)
    try {
      const response = await reportService.getStudentReport(studentId)
      if (response.status === 'success') {
        setStudentReport(response.data)
      }
    } catch (error) {
      console.error('Error fetching student report:', error)
      alert('Failed to load student report')
    } finally {
      setLoadingReport(false)
    }
  }

  const handleStudentSelect = (student) => {
    setSelectedStudent(student)
    fetchStudentReport(student.user_id || student.id)
  }

  const filteredStudents = students.filter(student =>
    student.user?.full_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
    student.student_id?.toLowerCase().includes(searchQuery.toLowerCase())
  )

  if (loading) {
    return (
      <DashboardLayout title="Reports & Analytics">
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
        </div>
      </DashboardLayout>
    )
  }

  return (
    <DashboardLayout title="Reports & Analytics">
      <div className="space-y-6">
        {/* Overview Stats */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div className="card">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600 mb-1">Total Students</p>
                <p className="text-3xl font-bold text-gray-900">{stats?.students || 0}</p>
              </div>
              <div className="bg-blue-500 p-3 rounded-lg">
                <GraduationCap className="w-8 h-8 text-white" />
              </div>
            </div>
          </div>

          <div className="card">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600 mb-1">Total Lecturers</p>
                <p className="text-3xl font-bold text-gray-900">{stats?.lecturers || 0}</p>
              </div>
              <div className="bg-green-500 p-3 rounded-lg">
                <Users className="w-8 h-8 text-white" />
              </div>
            </div>
          </div>

          <div className="card">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600 mb-1">Total Achievements</p>
                <p className="text-3xl font-bold text-gray-900">{stats?.achievements ? Object.values(stats.achievements).reduce((a, b) => a + b, 0) : 0}</p>
              </div>
              <div className="bg-purple-500 p-3 rounded-lg">
                <Award className="w-8 h-8 text-white" />
              </div>
            </div>
          </div>

          <div className="card">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600 mb-1">Pending Review</p>
                <p className="text-3xl font-bold text-gray-900">{stats?.achievements?.pending_verification || 0}</p>
              </div>
              <div className="bg-yellow-500 p-3 rounded-lg">
                <Clock className="w-8 h-8 text-white" />
              </div>
            </div>
          </div>
        </div>

        {/* Achievement Status Breakdown */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="card">
            <h2 className="text-xl font-bold text-gray-900 mb-4 flex items-center">
              <BarChart3 className="w-5 h-5 mr-2 text-primary-600" />
              Achievement by Status
            </h2>
            <div className="space-y-4">
              {stats?.achievements && Object.entries(stats.achievements).map(([status, count]) => {
                const total = Object.values(stats.achievements).reduce((a, b) => a + b, 0)
                const percentage = total > 0 ? (count / total * 100).toFixed(1) : 0
                const colors = {
                  draft: 'bg-blue-500',
                  pending_verification: 'bg-yellow-500',
                  verified: 'bg-green-500',
                  rejected: 'bg-red-500'
                }
                return (
                  <div key={status}>
                    <div className="flex justify-between text-sm mb-1">
                      <span className="text-gray-600 capitalize">{status.replace('_', ' ')}</span>
                      <span className="font-medium">{count} ({percentage}%)</span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2.5">
                      <div
                        className={`h-2.5 rounded-full ${colors[status] || 'bg-gray-500'}`}
                        style={{ width: `${percentage}%` }}
                      ></div>
                    </div>
                  </div>
                )
              })}
            </div>
          </div>

          <div className="card">
            <h2 className="text-xl font-bold text-gray-900 mb-4 flex items-center">
              <Trophy className="w-5 h-5 mr-2 text-primary-600" />
              Achievement Summary
            </h2>
            <div className="grid grid-cols-2 gap-4">
              <div className="p-4 bg-green-50 rounded-lg text-center">
                <p className="text-2xl font-bold text-green-600">{stats?.achievements?.verified || 0}</p>
                <p className="text-sm text-gray-600">Verified</p>
              </div>
              <div className="p-4 bg-yellow-50 rounded-lg text-center">
                <p className="text-2xl font-bold text-yellow-600">{stats?.achievements?.pending_verification || 0}</p>
                <p className="text-sm text-gray-600">Pending</p>
              </div>
              <div className="p-4 bg-blue-50 rounded-lg text-center">
                <p className="text-2xl font-bold text-blue-600">{stats?.achievements?.draft || 0}</p>
                <p className="text-sm text-gray-600">Draft</p>
              </div>
              <div className="p-4 bg-red-50 rounded-lg text-center">
                <p className="text-2xl font-bold text-red-600">{stats?.achievements?.rejected || 0}</p>
                <p className="text-sm text-gray-600">Rejected</p>
              </div>
            </div>
          </div>
        </div>

        {/* Top Performers */}
        {stats?.top_performers && stats.top_performers.length > 0 && (
          <div className="card">
            <h2 className="text-xl font-bold text-gray-900 mb-4 flex items-center">
              <TrendingUp className="w-5 h-5 mr-2 text-primary-600" />
              Top Performers
            </h2>
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Rank</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Student</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Program</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Total</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Verified</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {stats.top_performers.map((performer, index) => (
                    <tr key={performer.student_id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`inline-flex items-center justify-center w-8 h-8 rounded-full ${index === 0 ? 'bg-yellow-100 text-yellow-800' :
                          index === 1 ? 'bg-gray-100 text-gray-800' :
                            index === 2 ? 'bg-orange-100 text-orange-800' :
                              'bg-gray-50 text-gray-600'
                          }`}>
                          {index + 1}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="font-medium text-gray-900">{performer.full_name}</div>
                        <div className="text-sm text-gray-500">{performer.student_id}</div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                        {performer.program_study}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {performer.total_achievements}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className="badge badge-success">{performer.verified_achievements}</span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {/* Student Report Search */}
        <div className="card">
          <h2 className="text-xl font-bold text-gray-900 mb-4 flex items-center">
            <FileText className="w-5 h-5 mr-2 text-primary-600" />
            Student Report
          </h2>

          <div className="relative mb-4">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
            <input
              type="text"
              placeholder="Search student by name or ID..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="input pl-10"
            />
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Student List */}
            <div className="border border-gray-200 rounded-lg max-h-96 overflow-y-auto">
              {filteredStudents.length > 0 ? (
                filteredStudents.map((student) => (
                  <div
                    key={student.id}
                    onClick={() => handleStudentSelect(student)}
                    className={`p-4 border-b border-gray-100 cursor-pointer hover:bg-gray-50 transition-colors ${selectedStudent?.id === student.id ? 'bg-primary-50 border-l-4 border-l-primary-500' : ''
                      }`}
                  >
                    <div className="font-medium text-gray-900">{student.user?.full_name}</div>
                    <div className="text-sm text-gray-500">{student.student_id} • {student.program_study}</div>
                  </div>
                ))
              ) : (
                <div className="p-8 text-center text-gray-500">
                  {searchQuery ? 'No students found' : 'No students available'}
                </div>
              )}
            </div>

            {/* Student Report Detail */}
            <div className="border border-gray-200 rounded-lg p-4">
              {loadingReport ? (
                <div className="flex items-center justify-center h-64">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
                </div>
              ) : studentReport ? (
                <div className="space-y-4">
                  <div className="border-b pb-4">
                    <h3 className="font-bold text-lg text-gray-900">
                      {studentReport.student?.full_name}
                    </h3>
                    <p className="text-sm text-gray-500">
                      {studentReport.student?.student_id} • {studentReport.student?.program_study}
                    </p>
                    {studentReport.student?.advisor && (
                      <p className="text-sm text-gray-500 mt-1">
                        Advisor: {studentReport.student.advisor.full_name}
                      </p>
                    )}
                  </div>

                  <div className="grid grid-cols-2 gap-3">
                    <div className="p-3 bg-blue-50 rounded-lg text-center">
                      <p className="text-2xl font-bold text-blue-600">{studentReport.summary?.total_achievements || 0}</p>
                      <p className="text-xs text-gray-600">Total</p>
                    </div>
                    <div className="p-3 bg-green-50 rounded-lg text-center">
                      <p className="text-2xl font-bold text-green-600">{studentReport.summary?.verified_achievements || 0}</p>
                      <p className="text-xs text-gray-600">Verified</p>
                    </div>
                    <div className="p-3 bg-yellow-50 rounded-lg text-center">
                      <p className="text-2xl font-bold text-yellow-600">{studentReport.summary?.pending_achievements || 0}</p>
                      <p className="text-xs text-gray-600">Pending</p>
                    </div>
                    <div className="p-3 bg-red-50 rounded-lg text-center">
                      <p className="text-2xl font-bold text-red-600">{studentReport.summary?.rejected_achievements || 0}</p>
                      <p className="text-xs text-gray-600">Rejected</p>
                    </div>
                  </div>

                  {studentReport.achievements_by_level && (
                    <div>
                      <h4 className="font-medium text-gray-700 mb-2">By Level</h4>
                      <div className="flex flex-wrap gap-2">
                        {Object.entries(studentReport.achievements_by_level).map(([level, count]) => (
                          <span key={level} className="badge badge-info">
                            {level}: {count}
                          </span>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              ) : (
                <div className="flex flex-col items-center justify-center h-64 text-gray-500">
                  <FileText className="w-12 h-12 text-gray-300 mb-2" />
                  <p>Select a student to view report</p>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </DashboardLayout>
  )
}
