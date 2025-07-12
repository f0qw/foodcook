<template>
  <div class="meal-records-page">
    <div class="page-header">
      <h1>用餐记录</h1>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        添加记录
      </el-button>
    </div>

    <!-- 搜索和筛选 -->
    <div class="search-section">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            @change="handleSearch"
            style="width: 100%"
          />
        </el-col>
        <el-col :span="4">
          <el-button @click="resetSearch">重置</el-button>
        </el-col>
      </el-row>
    </div>

    <!-- 用餐记录列表 -->
    <div class="records-grid" v-loading="mealRecordsStore.loading">
      <el-card
        v-for="record in mealRecordsStore.mealRecords"
        :key="record.id"
        class="record-card"
      >
        <div class="record-header">
          <div class="record-date">
            <el-icon><Calendar /></el-icon>
            {{ formatDate(record.created_at) }}
          </div>
          <el-button
            size="small"
            type="danger"
            @click="deleteRecord(record.id)"
          >
            删除
          </el-button>
        </div>

        <div class="record-image" v-if="record.image_url">
          <el-image
            :src="record.image_url"
            fit="cover"
            style="width: 100%; height: 200px;"
            :preview-src-list="[record.image_url]"
          />
        </div>

        <div class="record-content">
          <div class="dishes-list">
            <h4>用餐菜品:</h4>
            <div class="dish-items">
              <el-tag
                v-for="dish in record.dishes"
                :key="dish.id"
                class="dish-tag"
              >
                {{ dish.name }} (¥{{ dish.price }})
              </el-tag>
            </div>
          </div>

          <div class="total-price">
            <strong>总价格: ¥{{ calculateTotal(record.dishes) }}</strong>
          </div>

          <div class="thoughts" v-if="record.thoughts">
            <h4>用餐感想:</h4>
            <p>{{ record.thoughts }}</p>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 分页 -->
    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="mealRecordsStore.currentPage"
        v-model:page-size="mealRecordsStore.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="mealRecordsStore.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 创建用餐记录对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      title="创建用餐记录"
      width="600px"
    >
      <el-form
        ref="recordFormRef"
        :model="recordForm"
        :rules="recordRules"
        label-width="100px"
      >
        <el-form-item label="选择菜品" prop="dish_ids">
          <el-checkbox-group v-model="recordForm.dish_ids">
            <el-checkbox
              v-for="dish in availableDishes"
              :key="dish.id"
              :label="dish.id"
            >
              {{ dish.name }} (¥{{ dish.price }})
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        
        <el-form-item label="用餐感想" prop="thoughts">
          <el-input
            v-model="recordForm.thoughts"
            type="textarea"
            :rows="3"
            placeholder="记录一下今天的用餐感受..."
          />
        </el-form-item>
        
        <el-form-item label="图片链接" prop="image_url">
          <el-input
            v-model="recordForm.image_url"
            placeholder="可选：添加用餐图片链接"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCreateDialog = false">取消</el-button>
          <el-button type="primary" @click="saveRecord" :loading="saving">
            创建记录
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useMealRecordsStore } from '@/stores/mealRecords'
import { useDishesStore } from '@/stores/dishes'
import { ElMessage, ElMessageBox } from 'element-plus'

const mealRecordsStore = useMealRecordsStore()
const dishesStore = useDishesStore()

// 响应式数据
const dateRange = ref([])
const showCreateDialog = ref(false)
const saving = ref(false)
const recordFormRef = ref()

const recordForm = reactive({
  dish_ids: [],
  thoughts: '',
  image_url: ''
})

const recordRules = {
  dish_ids: [
    { required: true, message: '请至少选择一个菜品', trigger: 'change' }
  ]
}

// 计算属性
const availableDishes = computed(() => dishesStore.dishes)

// 方法
const handleSearch = () => {
  const params = {}
  if (dateRange.value && dateRange.value.length === 2) {
    params.start_date = dateRange.value[0]
    params.end_date = dateRange.value[1]
  }
  mealRecordsStore.getMealRecords(params)
}

const resetSearch = () => {
  dateRange.value = []
  mealRecordsStore.getMealRecords()
}

const handleSizeChange = (size) => {
  mealRecordsStore.setPageSize(size)
  mealRecordsStore.getMealRecords()
}

const handleCurrentChange = (page) => {
  mealRecordsStore.setPage(page)
  mealRecordsStore.getMealRecords()
}

const saveRecord = async () => {
  if (!recordFormRef.value) return

  try {
    await recordFormRef.value.validate()
    saving.value = true

    await mealRecordsStore.createMealRecord(recordForm)
    showCreateDialog.value = false
    resetForm()
  } catch (error) {
    console.error('创建用餐记录失败:', error)
  } finally {
    saving.value = false
  }
}

const deleteRecord = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除这个用餐记录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await mealRecordsStore.deleteMealRecord(id)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除用餐记录失败:', error)
    }
  }
}

const resetForm = () => {
  Object.assign(recordForm, {
    dish_ids: [],
    thoughts: '',
    image_url: ''
  })
  if (recordFormRef.value) {
    recordFormRef.value.resetFields()
  }
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const calculateTotal = (dishes) => {
  return dishes.reduce((sum, dish) => sum + dish.price, 0).toFixed(2)
}

// 生命周期
onMounted(() => {
  mealRecordsStore.getMealRecords()
  dishesStore.getDishes()
})
</script>

<style scoped>
.meal-records-page {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h1 {
  margin: 0;
  color: #333;
}

.search-section {
  margin-bottom: 20px;
  padding: 20px;
  background-color: #f5f5f5;
  border-radius: 8px;
}

.records-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.record-card {
  transition: transform 0.3s, box-shadow 0.3s;
}

.record-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.record-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.record-date {
  display: flex;
  align-items: center;
  gap: 5px;
  color: #666;
  font-size: 14px;
}

.record-image {
  margin-bottom: 15px;
  border-radius: 8px;
  overflow: hidden;
}

.record-content {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.dishes-list h4 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 16px;
}

.dish-items {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.dish-tag {
  margin: 0;
}

.total-price {
  padding: 10px;
  background-color: #f0f9ff;
  border-radius: 6px;
  text-align: center;
  color: #f56c6c;
  font-size: 16px;
}

.thoughts h4 {
  margin: 0 0 8px 0;
  color: #333;
  font-size: 16px;
}

.thoughts p {
  margin: 0;
  color: #666;
  line-height: 1.6;
  background-color: #f9f9f9;
  padding: 10px;
  border-radius: 6px;
}

.pagination-wrapper {
  text-align: center;
  margin-top: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 