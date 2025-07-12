import { defineStore } from 'pinia'
import { ref } from 'vue'
import { mealRecordsAPI } from '@/utils/api'
import { ElMessage } from 'element-plus'

export const useMealRecordsStore = defineStore('mealRecords', () => {
  const mealRecords = ref([])
  const loading = ref(false)
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(10)

  // 获取用餐记录列表
  const getMealRecords = async (params = {}) => {
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
      await getMealRecords()
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
      await getMealRecords()
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
      await getMealRecords()
    } catch (error) {
      throw error
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
    mealRecords,
    loading,
    total,
    currentPage,
    pageSize,
    getMealRecords,
    getMealRecordById,
    createMealRecord,
    updateMealRecord,
    deleteMealRecord,
    setPage,
    setPageSize
  }
}) 