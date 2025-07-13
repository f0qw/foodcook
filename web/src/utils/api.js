import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建axios实例
const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    // 添加token到请求头
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

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    const { response } = error
    
    if (response) {
      const { status, data } = response
      
      switch (status) {
        case 401:
          ElMessage.error('登录已过期，请重新登录')
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          window.location.href = '/login'
          break
        case 403:
          ElMessage.error('没有权限访问')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          ElMessage.error(data?.error || '请求失败')
      }
    } else {
      ElMessage.error('网络连接失败')
    }
    
    return Promise.reject(error)
  }
)

// API方法
export const authAPI = {
  // 用户注册
  register: (data) => api.post('/auth/register', data),
  
  // 用户登录
  login: (data) => api.post('/auth/login', data),
  
  // 获取用户信息
  getProfile: () => api.get('/auth/profile')
}

export const dishesAPI = {
  // 获取菜品列表
  getList: (params) => api.get('/dishes', { params }),
  
  // 获取菜品详情
  getById: (id) => api.get(`/dishes/${id}`),
  
  // 创建菜品
  create: (data) => api.post('/dishes', data),
  
  // 更新菜品
  update: (id, data) => api.put(`/dishes/${id}`, data),
  
  // 删除菜品
  delete: (id) => api.delete(`/dishes/${id}`),
  
  // 搜索菜品
  search: (params) => api.get('/dishes/search', { params })
}

export const ingredientsAPI = {
  // 获取食材列表
  getList: (params) => api.get('/ingredients', { params }),
  
  // 获取食材详情
  getById: (id) => api.get(`/ingredients/${id}`),
  
  // 创建食材
  create: (data) => api.post('/ingredients', data),
  
  // 更新食材
  update: (id, data) => api.put(`/ingredients/${id}`, data),
  
  // 删除食材
  delete: (id) => api.delete(`/ingredients/${id}`),
  
  // 搜索食材
  search: (params) => api.get('/ingredients/search', { params })
}

export const categoriesAPI = {
  // 获取分类列表
  getList: () => api.get('/categories'),
  
  // 获取分类详情
  getById: (id) => api.get(`/categories/${id}`),
  
  // 创建分类
  create: (data) => api.post('/categories', data),
  
  // 更新分类
  update: (id, data) => api.put(`/categories/${id}`, data),
  
  // 删除分类
  delete: (id) => api.delete(`/categories/${id}`)
}

export const mealRecordsAPI = {
  // 获取用餐记录列表
  getList: (params) => api.get('/meal-records', { params }),
  
  // 获取用餐记录详情
  getById: (id) => api.get(`/meal-records/${id}`),
  
  // 创建用餐记录
  create: (data) => api.post('/meal-records', data),
  
  // 更新用餐记录
  update: (id, data) => api.put(`/meal-records/${id}`, data),
  
  // 删除用餐记录
  delete: (id) => api.delete(`/meal-records/${id}`)
}

export default api 