import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8082/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api

// Auth APIs
export const authApi = {
  register: (data) => api.post('/auth/register', data),
  login: (data) => api.post('/auth/login', data),
  getCurrentUser: () => api.get('/me')
}

// Task APIs
export const taskApi = {
  getAll: (params) => api.get('/tasks', { params }),
  getById: (id) => api.get(`/tasks/${id}`),
  create: (data) => api.post('/tasks', data),
  getMy: () => api.get('/tasks/my'),
  update: (id, data) => api.put(`/tasks/${id}`, data),
  delete: (id) => api.delete(`/tasks/${id}`),
  complete: (id) => api.post(`/tasks/${id}/complete`)
}

// Bid APIs
export const bidApi = {
  create: (taskId, data) => api.post(`/bids/task/${taskId}`, data),
  getByTask: (taskId) => api.get(`/bids/task/${taskId}`),
  getMy: () => api.get('/bids/my'),
  accept: (bidId) => api.post(`/bids/${bidId}/accept`)
}

// Review APIs
export const reviewApi = {
  create: (taskId, data) => api.post(`/reviews/task/${taskId}`, data),
  getByEmployee: (employeeId) => api.get(`/reviews/employee/${employeeId}`),
  getMy: () => api.get('/reviews/my')
}

// Favorite APIs
export const favoriteApi = {
  toggle: (taskId) => api.post(`/favorites/task/${taskId}`),
  getMy: () => api.get('/favorites/my'),
  check: (taskId) => api.get(`/favorites/task/${taskId}/check`)
}

// Admin APIs
export const adminApi = {
  getUsers: (params) => api.get('/admin/users', { params }),
  getUserById: (id) => api.get(`/admin/users/${id}`),
  updateUser: (id, data) => api.put(`/admin/users/${id}`, data),
  deleteUser: (id) => api.delete(`/admin/users/${id}`),
  getTasks: (params) => api.get('/admin/tasks', { params }),
  updateTaskStatus: (id, data) => api.put(`/admin/tasks/${id}/status`, data),
  deleteTask: (id) => api.delete(`/admin/tasks/${id}`)
}
