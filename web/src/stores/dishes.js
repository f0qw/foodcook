import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { dishesAPI } from '@/utils/api'
import { ElMessage } from 'element-plus'

export const useDishesStore = defineStore('dishes', () => {
  const dishes = ref([])
  const loading = ref(false)
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(10)

  // 计算属性
  const hasMore = computed(() => {
    return dishes.value.length < total.value
  })

  // 获取菜品列表
  const fetchDishes = async (params = {}) => {
    loading.value = true
    try {
      const response = await dishesAPI.getList({
        offset: (currentPage.value - 1) * pageSize.value,
        limit: pageSize.value,
        ...params
      })
      dishes.value = response.data
      total.value = response.total
      return response
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取菜品列表（兼容旧方法）
  const getDishes = async (params = {}) => {
    return await fetchDishes(params)
  }

  // 加载更多菜品
  const loadMore = async () => {
    if (loading.value || !hasMore.value) return
    
    loading.value = true
    try {
      const nextPage = currentPage.value + 1
      const response = await dishesAPI.getList({
        offset: (nextPage - 1) * pageSize.value,
        limit: pageSize.value
      })
      
      dishes.value = [...dishes.value, ...response.data]
      currentPage.value = nextPage
      total.value = response.total
      return response
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取菜品详情
  const getDishById = async (id) => {
    try {
      const response = await dishesAPI.getById(id)
      return response
    } catch (error) {
      throw error
    }
  }

  // 创建菜品
  const createDish = async (data) => {
    try {
      const response = await dishesAPI.create(data)
      ElMessage.success('菜品创建成功')
      await fetchDishes()
      return response
    } catch (error) {
      throw error
    }
  }

  // 更新菜品
  const updateDish = async (id, data) => {
    try {
      const response = await dishesAPI.update(id, data)
      ElMessage.success('菜品更新成功')
      await fetchDishes()
      return response
    } catch (error) {
      throw error
    }
  }

  // 删除菜品
  const deleteDish = async (id) => {
    try {
      await dishesAPI.delete(id)
      ElMessage.success('菜品删除成功')
      await fetchDishes()
    } catch (error) {
      throw error
    }
  }

  // 搜索菜品
  const searchDishes = async (keyword) => {
    loading.value = true
    try {
      const response = await dishesAPI.search({
        q: keyword,
        offset: 0,
        limit: 50
      })
      dishes.value = response.data
      total.value = response.total
      return response
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // 设置分页
  const setCurrentPage = (page) => {
    currentPage.value = page
  }

  const setPageSize = (size) => {
    pageSize.value = size
    currentPage.value = 1
  }

  return {
    dishes,
    loading,
    total,
    currentPage,
    pageSize,
    hasMore,
    fetchDishes,
    getDishes,
    loadMore,
    getDishById,
    createDish,
    updateDish,
    deleteDish,
    searchDishes,
    setCurrentPage,
    setPageSize
  }
}) 