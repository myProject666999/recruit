import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '../services/api'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  const isLoggedIn = computed(() => !!token.value)
  const isEmployer = computed(() => user.value?.role === 'employer')
  const isEmployee = computed(() => user.value?.role === 'employee')
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function login(credentials) {
    const response = await authApi.login(credentials)
    const { token: newToken, user: newUser } = response.data
    
    token.value = newToken
    user.value = newUser
    
    localStorage.setItem('token', newToken)
    localStorage.setItem('user', JSON.stringify(newUser))
    
    return response.data
  }

  async function register(userData) {
    const response = await authApi.register(userData)
    const { token: newToken, user: newUser } = response.data
    
    token.value = newToken
    user.value = newUser
    
    localStorage.setItem('token', newToken)
    localStorage.setItem('user', JSON.stringify(newUser))
    
    return response.data
  }

  async function fetchCurrentUser() {
    if (!token.value) return null
    
    try {
      const response = await authApi.getCurrentUser()
      user.value = response.data
      localStorage.setItem('user', JSON.stringify(response.data))
      return response.data
    } catch (error) {
      logout()
      throw error
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return {
    token,
    user,
    isLoggedIn,
    isEmployer,
    isEmployee,
    isAdmin,
    login,
    register,
    fetchCurrentUser,
    logout
  }
})
