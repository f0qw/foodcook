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
            :src="record.image_url || '/default-dish.jpg'"
            fit="cover"
            style="width: 100%; height: 200px;"
            :preview-src-list="[record.image_url || '/default-dish.jpg']"
          >
            <template #error>
              <div class="image-placeholder">
                <el-icon><Picture /></el-icon>
                <span>暂无图片</span>
              </div>
            </template>
          </el-image>
        </div>

        <div class="record-content">
          <div class="dishes-list">
            <h4>用餐菜品:</h4>
            <div class="dish-items">
              <el-tag
                v-for="dishItem in record.dishes"
                :key="dishItem.id"
                class="dish-tag"
                @click="showDishDetail(dishItem.dish)"
                style="cursor: pointer;"
              >
                {{ dishItem.dish.name }} (¥{{ dishItem.dish.price }})
              </el-tag>
            </div>
          </div>

          <div class="total-price">
            <strong>总价格: ¥{{ record.total_price }}</strong>
          </div>

          <div class="thoughts" v-if="record.thoughts">
            <h4>用餐感想:</h4>
            <p>{{ record.thoughts }}</p>
          </div>

          <div class="record-actions">
            <el-button 
              size="small" 
              type="primary" 
              @click="showRecordDetail(record)"
            >
              查看详情
            </el-button>
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

    <!-- 用餐记录详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="用餐记录详情"
      width="800px"
    >
      <div v-if="selectedRecord" class="record-detail">
        <div class="detail-header">
          <div class="detail-date">
            <el-icon><Calendar /></el-icon>
            {{ formatDate(selectedRecord.created_at) }}
          </div>
          <div class="detail-price">
            <strong>总价格: ¥{{ selectedRecord.total_price }}</strong>
          </div>
        </div>

        <div class="detail-dishes">
          <h3>用餐菜品详情</h3>
          <div class="dishes-grid">
            <div 
              v-for="dishItem in selectedRecord.dishes" 
              :key="dishItem.id"
              class="dish-card"
            >
              <div class="dish-image">
                <el-image
                  :src="dishItem.dish.image_url || '/default-dish.jpg'"
                  fit="cover"
                  style="width: 100%; height: 150px;"
                  :preview-src-list="[dishItem.dish.image_url || '/default-dish.jpg']"
                >
                  <template #error>
                    <div class="image-placeholder">
                      <el-icon><Picture /></el-icon>
                      <span>暂无图片</span>
                    </div>
                  </template>
                </el-image>
              </div>
              <div class="dish-info">
                <h4>{{ dishItem.dish.name }}</h4>
                <p class="dish-description">{{ dishItem.dish.description }}</p>
                <div class="dish-meta">
                  <span class="dish-price">¥{{ dishItem.dish.price }}</span>
                  <span class="dish-quantity">数量: {{ dishItem.quantity }}</span>
                </div>
                <div class="dish-category" v-if="dishItem.dish.category">
                  <el-tag size="small">{{ dishItem.dish.category.name }}</el-tag>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="detail-thoughts" v-if="selectedRecord.thoughts">
          <h3>用餐感想</h3>
          <p>{{ selectedRecord.thoughts }}</p>
        </div>

        <div class="detail-image" v-if="selectedRecord.image_url">
          <h3>用餐图片</h3>
          <el-image
            :src="selectedRecord.image_url"
            fit="cover"
            style="width: 100%; max-height: 300px;"
            :preview-src-list="[selectedRecord.image_url]"
          />
        </div>
      </div>
    </el-dialog>

    <!-- 菜品详情对话框 -->
    <el-dialog
      v-model="showDishDetailDialog"
      title="菜品详情"
      width="600px"
    >
      <div v-if="selectedDish" class="dish-detail">
        <div class="dish-detail-image">
          <el-image
            :src="selectedDish.image_url || '/default-dish.jpg'"
            fit="cover"
            style="width: 100%; height: 250px;"
            :preview-src-list="[selectedDish.image_url || '/default-dish.jpg']"
          >
            <template #error>
              <div class="image-placeholder">
                <el-icon><Picture /></el-icon>
                <span>暂无图片</span>
              </div>
            </template>
          </el-image>
        </div>
        <div class="dish-detail-info">
          <h2>{{ selectedDish.name }}</h2>
          <p class="dish-description">{{ selectedDish.description }}</p>
          <div class="dish-detail-meta">
            <div class="meta-item">
              <strong>价格:</strong> ¥{{ selectedDish.price }}
            </div>
            <div class="meta-item" v-if="selectedDish.category">
              <strong>分类:</strong> {{ selectedDish.category.name }}
            </div>
            <div class="meta-item" v-if="selectedDish.cooking_link">
              <strong>烹饪链接:</strong> 
              <a :href="selectedDish.cooking_link" target="_blank">查看食谱</a>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useMealRecordsStore } from '@/stores/mealRecords'
import { useDishesStore } from '@/stores/dishes'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Calendar, Picture } from '@element-plus/icons-vue'

const mealRecordsStore = useMealRecordsStore()
const dishesStore = useDishesStore()

// 响应式数据
const dateRange = ref([])
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const showDishDetailDialog = ref(false)
const selectedRecord = ref(null)
const selectedDish = ref(null)
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

const showRecordDetail = (record) => {
  selectedRecord.value = record
  showDetailDialog.value = true
}

const showDishDetail = (dish) => {
  selectedDish.value = dish
  showDishDetailDialog.value = true
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

/* 记录详情样式 */
.record-detail {
  padding: 20px 0;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #eee;
}

.detail-date {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #666;
  font-size: 16px;
}

.detail-price {
  color: #f56c6c;
  font-size: 18px;
}

.detail-dishes h3 {
  margin: 0 0 15px 0;
  color: #333;
}

.dishes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 15px;
  margin-bottom: 20px;
}

.dish-card {
  border: 1px solid #eee;
  border-radius: 8px;
  overflow: hidden;
  transition: box-shadow 0.3s;
}

.dish-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.dish-image {
  position: relative;
}

.image-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 150px;
  background-color: #f5f5f5;
  color: #999;
  font-size: 14px;
}

.dish-info {
  padding: 15px;
}

.dish-info h4 {
  margin: 0 0 8px 0;
  color: #333;
  font-size: 16px;
}

.dish-description {
  margin: 0 0 10px 0;
  color: #666;
  font-size: 14px;
  line-height: 1.4;
}

.dish-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.dish-price {
  color: #f56c6c;
  font-weight: bold;
}

.dish-quantity {
  color: #666;
  font-size: 14px;
}

.dish-category {
  margin-top: 8px;
}

.detail-thoughts {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.detail-thoughts h3 {
  margin: 0 0 10px 0;
  color: #333;
}

.detail-thoughts p {
  margin: 0;
  color: #666;
  line-height: 1.6;
  background-color: #f9f9f9;
  padding: 15px;
  border-radius: 6px;
}

.detail-image {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.detail-image h3 {
  margin: 0 0 15px 0;
  color: #333;
}

/* 菜品详情样式 */
.dish-detail {
  padding: 20px 0;
}

.dish-detail-image {
  margin-bottom: 20px;
  border-radius: 8px;
  overflow: hidden;
}

.dish-detail-info h2 {
  margin: 0 0 15px 0;
  color: #333;
}

.dish-detail-meta {
  margin-top: 20px;
}

.meta-item {
  margin-bottom: 10px;
  color: #666;
}

.meta-item strong {
  color: #333;
  margin-right: 8px;
}

.meta-item a {
  color: #409eff;
  text-decoration: none;
}

.meta-item a:hover {
  text-decoration: underline;
}

.record-actions {
  margin-top: 15px;
  text-align: center;
}
</style> 