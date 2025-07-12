import { defineStore } from 'pinia'
import { ref } from 'vue'
import { ingredientsAPI } from '@/utils/api'
import { ElMessage } from 'element-plus'

export const useIngredientsStore = defineStore('ingredients', () => {
  const ingredients = ref([])
  const loading = ref(false)
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(10)

  // 获取食材列表
  const getIngredients = async (params = {}) => {
    loading.value = true
    try {
      const response = await ingredientsAPI.getList({
        offset: (currentPage.value - 1) * pageSize.value,
        limit: pageSize.value,
        ...params
      })
      ingredients.value = response.data
      total.value = response.total
      return response
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取食材详情
  const getIngredientById = async (id) => {
    try {
      const response = await ingredientsAPI.getById(id)
      return response
    } catch (error) {
      throw error
    }
  }

  // 创建食材
  const createIngredient = async (data) => {
    try {
      const response = await ingredientsAPI.create(data)
      ElMessage.success('食材创建成功')
      await getIngredients()
      return response
    } catch (error) {
      throw error
    }
  }

  // 更新食材
  const updateIngredient = async (id, data) => {
    try {
      const response = await ingredientsAPI.update(id, data)
      ElMessage.success('食材更新成功')
      await getIngredients()
      return response
    } catch (error) {
      throw error
    }
  }

  // 删除食材
  const deleteIngredient = async (id) => {
    try {
      await ingredientsAPI.delete(id)
      ElMessage.success('食材删除成功')
      await getIngredients()
    } catch (error) {
      throw error
    }
  }

  // 搜索食材
  const searchIngredients = async (keyword) => {
    loading.value = true
    try {
      const response = await ingredientsAPI.search({
        q: keyword,
        offset: 0,
        limit: 50
      })
      ingredients.value = response.data
      total.value = response.total
      return response
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // 设置分页
  const setPage = (page) => {
    currentPage.value = page
  }

  const setPageSize = (size) => {
    pageSize.value = size
    currentPage.value = 1
  }

  return {
    ingredients,
    loading,
    total,
    currentPage,
    pageSize,
    getIngredients,
    getIngredientById,
    createIngredient,
    updateIngredient,
    deleteIngredient,
    searchIngredients,
    setPage,
    setPageSize
  }
}) 