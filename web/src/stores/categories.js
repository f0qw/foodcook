import { defineStore } from 'pinia'
import { categoriesAPI } from '@/utils/api'

export const useCategoriesStore = defineStore('categories', {
  state: () => ({
    categories: [],
    loading: false,
    error: null
  }),

  getters: {
    getCategoryById: (state) => (id) => {
      return state.categories.find(category => category.id === id)
    }
  },

  actions: {
    // 获取分类列表
    async getCategories() {
      this.loading = true
      this.error = null
      
      try {
        const response = await categoriesAPI.getList()
        this.categories = response.data.data || response.data
      } catch (error) {
        console.error('获取分类列表失败:', error)
        this.error = error.response?.data?.error || '获取分类列表失败'
        throw error
      } finally {
        this.loading = false
      }
    },

    // 创建分类
    async createCategory(categoryData) {
      this.loading = true
      this.error = null
      
      try {
        const response = await categoriesAPI.create(categoryData)
        // 重新获取分类列表
        await this.getCategories()
        return response.data
      } catch (error) {
        console.error('创建分类失败:', error)
        this.error = error.response?.data?.error || '创建分类失败'
        throw error
      } finally {
        this.loading = false
      }
    },

    // 更新分类
    async updateCategory(id, categoryData) {
      this.loading = true
      this.error = null
      
      try {
        const response = await categoriesAPI.update(id, categoryData)
        // 重新获取分类列表
        await this.getCategories()
        return response.data
      } catch (error) {
        console.error('更新分类失败:', error)
        this.error = error.response?.data?.error || '更新分类失败'
        throw error
      } finally {
        this.loading = false
      }
    },

    // 删除分类
    async deleteCategory(id) {
      this.loading = true
      this.error = null
      
      try {
        await categoriesAPI.delete(id)
        // 重新获取分类列表
        await this.getCategories()
      } catch (error) {
        console.error('删除分类失败:', error)
        this.error = error.response?.data?.error || '删除分类失败'
        throw error
      } finally {
        this.loading = false
      }
    }
  }
}) 