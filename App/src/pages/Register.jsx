import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { GraduationCap, AlertCircle, ArrowLeft } from 'lucide-react'

export default function Register() {
  const { user } = useAuth()
  const navigate = useNavigate()

  // Redirect if already logged in
  useEffect(() => {
    if (user) {
      switch (user.role_name) {
        case 'Admin':
          navigate('/admin/dashboard', { replace: true })
          break
        case 'Mahasiswa':
          navigate('/student/dashboard', { replace: true })
          break
        case 'Dosen Wali':
          navigate('/lecturer/dashboard', { replace: true })
          break
        default:
          break
      }
    }
  }, [user, navigate])

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100 px-4 py-12">
      <div className="max-w-md w-full">
        <div className="text-center mb-8">
          <div className="flex justify-center mb-4">
            <div className="bg-primary-600 p-3 rounded-full">
              <GraduationCap className="w-12 h-12 text-white" />
            </div>
          </div>
          <h1 className="text-3xl font-bold text-gray-900">Student Achievement System</h1>
          <p className="mt-2 text-gray-600">Account Registration</p>
        </div>

        <div className="card">
          {/* Information Notice */}
          <div className="p-4 bg-blue-50 border border-blue-200 rounded-lg mb-6">
            <div className="flex items-start space-x-3">
              <AlertCircle className="w-6 h-6 text-blue-500 flex-shrink-0 mt-0.5" />
              <div>
                <h3 className="font-semibold text-blue-800">Registration Information</h3>
                <p className="text-sm text-blue-700 mt-1">
                  User accounts are managed by the System Administrator. If you need an account, please contact your administrator or academic office.
                </p>
              </div>
            </div>
          </div>

          {/* Account Types */}
          <div className="space-y-4 mb-6">
            <h3 className="font-semibold text-gray-900">Available Account Types:</h3>

            <div className="p-4 border border-gray-200 rounded-lg">
              <h4 className="font-medium text-gray-900">👨‍🎓 Mahasiswa (Student)</h4>
              <p className="text-sm text-gray-600 mt-1">
                For students to submit and track their academic achievements
              </p>
            </div>

            <div className="p-4 border border-gray-200 rounded-lg">
              <h4 className="font-medium text-gray-900">👨‍🏫 Dosen Wali (Academic Advisor)</h4>
              <p className="text-sm text-gray-600 mt-1">
                For lecturers to verify and manage their advisees' achievements
              </p>
            </div>

            <div className="p-4 border border-gray-200 rounded-lg">
              <h4 className="font-medium text-gray-900">👨‍💼 Admin</h4>
              <p className="text-sm text-gray-600 mt-1">
                For system administrators to manage users and reports
              </p>
            </div>
          </div>

          {/* Contact Info */}
          <div className="p-4 bg-yellow-50 border border-yellow-200 rounded-lg mb-6">
            <h3 className="font-semibold text-yellow-800">How to Get an Account:</h3>
            <ol className="text-sm text-yellow-700 mt-2 space-y-1 list-decimal list-inside">
              <li>Contact the academic administration office</li>
              <li>Provide your student/staff ID and email</li>
              <li>Wait for account activation notification</li>
              <li>Use your credentials to login</li>
            </ol>
          </div>

          {/* Demo Accounts */}
          <div className="p-4 bg-gray-50 rounded-lg mb-6">
            <p className="text-sm font-medium text-gray-700 mb-2">Demo Accounts:</p>
            <div className="space-y-1 text-xs text-gray-600">
              <p><strong>Admin:</strong> admin / admin123</p>
              <p><strong>Student:</strong> student001 / student123</p>
              <p><strong>Lecturer:</strong> lecturer001 / lecturer123</p>
            </div>
          </div>

          {/* Back to Login */}
          <Link
            to="/login"
            className="btn btn-primary w-full flex items-center justify-center space-x-2"
          >
            <ArrowLeft className="w-4 h-4" />
            <span>Back to Login</span>
          </Link>
        </div>
      </div>
    </div>
  )
}
