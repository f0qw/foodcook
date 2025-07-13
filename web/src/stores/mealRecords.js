import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { mealRecordsAPI } from '@/utils/api'
import { ElMessage } from 'element-plus'

export const useMealRecordsStore = defineStore('mealRecords', () => {
  const mealRecords = ref([])
  const loading = ref(false)
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(10)

  // 计算属性
  const hasMore = computed(() => {
    return mealRecords.value.length < total.value
  })

  // 获取用餐记录列表
  const fetchMealRecords = async (params = {}) => {
    loading.value = true
    try {
      const response = await mealRecordsAPI.getList({
        offset: (currentPage.value - 1) * pageSize.value,
        limit: pageSize.value,
        ...params
      })
      mealRecords.value = response.data
      total.value = response.total
      return response
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取用餐记录列表（兼容旧方法）
  const getMealRecords = async (params = {}) => {
    return await fetchMealRecords(params)
  }

  // 加载更多用餐记录
  const loadMore = async () => {
    if (loading.value || !hasMore.value) return
    
    loading.value = true
    try {
      const nextPage = currentPage.value + 1
      const response = await mealRecordsAPI.getList({
        offset: (nextPage - 1) * pageSize.value,
        limit: pageSize.value
      })
      
      mealRecords.value = [...mealRecords.value, ...response.data]
      currentPage.value = nextPage
      total.value = response.total
      return response
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取用餐记录详情
  const getMealRecordById = async (id) => {
    try {
      const response = await mealRecordsAPI.getById(id)
      return response
    } catch (error) {
      throw error
    }
  }

  // 创建用餐记录
  const createMealRecord = async (data) => {
    try {
      const response = await mealRecordsAPI.create(data)
      ElMessage.success('用餐记录创建成功')
      await fetchMealRecords()
      return response
    } catch (error) {
      throw error
    }
  }

  // 更新用餐记录
  const updateMealRecord = async (id, data) => {
    try {
      const response = await mealRecordsAPI.update(id, data)
      ElMessage.success('用餐记录更新成功')
      await fetchMealRecords()
      return response
    } catch (error) {
      throw error
    }
  }

  // 删除用餐记录
  const deleteMealRecord = async (id) => {
    try {
      await mealRecordsAPI.delete(id)
      ElMessage.success('用餐记录删除成功')
      await fetchMealRecords()
    } catch (error) {
      throw error
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
    mealRecords,
    loading,
    total,
    currentPage,
    pageSize,
    hasMore,
    fetchMealRecords,
    getMealRecords,
    loadMore,
    getMealRecordById,
    createMealRecord,
    updateMealRecord,
    deleteMealRecord,
    setCurrentPage,
    setPageSize
  }
}) 