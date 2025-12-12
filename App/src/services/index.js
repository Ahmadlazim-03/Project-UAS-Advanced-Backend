import api from './api'

export const authService = {
  login: async (username, password) => {
    const response = await api.post('/auth/login', { username, password })
    return response.data
  },

  register: async (userData) => {
    const response = await api.post('/auth/register', userData)
    return response.data
  },

  logout: async () => {
    const response = await api.post('/auth/logout')
    return response.data
  },

  getProfile: async () => {
    const response = await api.get('/auth/profile')
    return response.data
  },

  refreshToken: async (refreshToken) => {
    const response = await api.post('/auth/refresh', { refresh_token: refreshToken })
    return response.data
  },
}

export const userService = {
  getUsers: async (page = 1, limit = 10) => {
    const response = await api.get(`/users?page=${page}&limit=${limit}`)
    return response.data
  },

  getUserById: async (id) => {
    const response = await api.get(`/users/${id}`)
    return response.data
  },

  createUser: async (userData) => {
    const response = await api.post('/users', userData)
    return response.data
  },

  updateUser: async (id, userData) => {
    const response = await api.put(`/users/${id}`, userData)
    return response.data
  },

  deleteUser: async (id) => {
    const response = await api.delete(`/users/${id}`)
    return response.data
  },

  assignRole: async (id, roleId) => {
    const response = await api.put(`/users/${id}/role`, { role_id: roleId })
    return response.data
  },
}

export const studentService = {
  getStudents: async (page = 1, limit = 10) => {
    const response = await api.get(`/students?page=${page}&limit=${limit}`)
    return response.data
  },

  getStudentById: async (id) => {
    const response = await api.get(`/students/${id}`)
    return response.data
  },

  getStudentAchievements: async (id) => {
    const response = await api.get(`/students/${id}/achievements`)
    return response.data
  },

  assignAdvisor: async (id, advisorId) => {
    const response = await api.put(`/students/${id}/advisor`, { advisor_id: advisorId })
    return response.data
  },
}

export const lecturerService = {
  getLecturers: async (page = 1, limit = 10) => {
    const response = await api.get(`/lecturers?page=${page}&limit=${limit}`)
    return response.data
  },

  getAdvisees: async (id) => {
    const response = await api.get(`/lecturers/${id}/advisees`)
    return response.data
  },

  getAdviseeAchievements: async () => {
    const response = await api.get('/lecturers/advisees/achievements')
    return response.data
  },
}

export const achievementService = {
  getAchievements: async (params = {}) => {
    const queryString = new URLSearchParams(params).toString()
    const response = await api.get(`/achievements?${queryString}`)
    return response.data
  },

  getAchievementById: async (id) => {
    const response = await api.get(`/achievements/${id}`)
    return response.data
  },

  createAchievement: async (achievementData) => {
    const response = await api.post('/achievements', achievementData)
    return response.data
  },

  updateAchievement: async (id, achievementData) => {
    const response = await api.put(`/achievements/${id}`, achievementData)
    return response.data
  },

  deleteAchievement: async (id) => {
    const response = await api.delete(`/achievements/${id}`)
    return response.data
  },

  submitAchievement: async (id) => {
    const response = await api.post(`/achievements/${id}/submit`)
    return response.data
  },

  verifyAchievement: async (id, comments) => {
    const response = await api.post(`/achievements/${id}/verify`, { comments })
    return response.data
  },

  rejectAchievement: async (id, reason) => {
    const response = await api.post(`/achievements/${id}/reject`, { reason })
    return response.data
  },
}

export const fileService = {
  uploadFile: async (file) => {
    const formData = new FormData()
    formData.append('file', file)
    
    const response = await api.post('/files/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response.data
  },

  deleteFile: async (filename) => {
    const response = await api.delete(`/files/${filename}`)
    return response.data
  },
}

export const reportService = {
  getStatistics: async () => {
    const response = await api.get('/reports/statistics')
    return response.data
  },

  getStudentReport: async (id) => {
    const response = await api.get(`/reports/student/${id}`)
    return response.data
  },
}
