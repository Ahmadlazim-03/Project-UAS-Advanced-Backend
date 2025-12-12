import { createContext, useContext, useState, useEffect } from 'react'
import { authService } from '../services'
import { setToken, setUser, getUser, getToken, clearAuth } from '../utils/auth'

const AuthContext = createContext(null)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}

export const AuthProvider = ({ children }) => {
  const [user, setUserState] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Check if user is logged in
    const token = getToken()
    const savedUser = getUser()

    if (token && savedUser) {
      setUserState(savedUser)
    }
    setLoading(false)
  }, [])

  const login = async (username, password) => {
    try {
      const response = await authService.login(username, password)
      if (response.status === 'success' && response.data) {
        setToken(response.data.token)
        // Map 'role' to 'role_name' since backend returns 'role' but components expect 'role_name'
        const userData = {
          ...response.data.user,
          role_name: response.data.user.role || response.data.user.role_name
        }
        setUser(userData)
        setUserState(userData)
        return { success: true }
      }
      return { success: false, message: response.message || 'Login failed' }
    } catch (error) {
      return {
        success: false,
        message: error.response?.data?.message || 'Login failed. Please try again.'
      }
    }
  }

  const logout = async () => {
    try {
      await authService.logout()
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      clearAuth()
      setUserState(null)
    }
  }

  const register = async (userData) => {
    try {
      const response = await authService.register(userData)
      if (response.status === 'success') {
        return { success: true }
      }
      return { success: false, message: response.message || 'Registration failed' }
    } catch (error) {
      return {
        success: false,
        message: error.response?.data?.message || 'Registration failed. Please try again.'
      }
    }
  }

  const value = {
    user,
    loading,
    login,
    logout,
    register,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
