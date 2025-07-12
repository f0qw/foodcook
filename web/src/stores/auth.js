import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authAPI } from '@/utils/api'
import { ElMessage } from 'element-plus'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || '')

  // 计算属性
  const isAuthenticated = computed(() => !!token.value)

  // 登录
  const login = async (credentials) => {
    try {
      const response = await authAPI.login(credentials)
      user.value = response.user
      token.value = response.token
      localStorage.setItem('token', response.token)
      localStorage.setItem('user', JSON.stringify(response.user))
      ElMessage.success('登录成功')
      return response
    } catch (error) {
      throw error
    }
  }

  // 注册
  const register = async (userData) => {
    try {
      const response = await authAPI.register(userData)
      user.value = response.user
      token.value = response.token
      localStorage.setItem('token', response.token)
      localStorage.setItem('user', JSON.stringify(response.user))
      ElMessage.success('注册成功')
      return response
    } catch (error) {
      throw error
    }
  }

  // 获取用户信息
  const getProfile = async () => {
    try {
      const response = await authAPI.getProfile()
      user.value = response
      localStorage.setItem('user', JSON.stringify(response))
      return response
    } catch (error) {
      throw error
    }
  }

  // 登出
  const logout = () => {
    user.value = null
    token.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    ElMessage.success('已退出登录')
  }

  // 初始化用户信息
  const initUser = () => {
    const savedUser = localStorage.getItem('user')
    if (savedUser) {
      user.value = JSON.parse(savedUser)
    }
  }

  return {
    user,
    token,
    isAuthenticated,
    login,
    register,
    getProfile,
    logout,
    initUser
  }
}) 