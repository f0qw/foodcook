<template>
  <div class="meal-records-page">
    <!-- 移动端头部 -->
    <div v-if="isMobile" class="mobile-header">
      <div class="mobile-header-content">
        <el-button @click="$router.go(-1)" type="text" class="back-btn">
          <el-icon><ArrowLeft /></el-icon>
        </el-button>
        <h2>用餐记录</h2>
        <el-button type="primary" size="small" @click="showCreateDialog = true">
          <el-icon><Plus /></el-icon>
        </el-button>
      </div>
    </div>

    <!-- 桌面端头部 -->
    <div v-else class="page-header">
      <h1>用餐记录</h1>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        添加记录
      </el-button>
    </div>

    <!-- 主要内容 -->
    <div class="main-content" :class="{ 'mobile-content': isMobile }">
      <!-- 搜索和筛选 -->
      <div class="search-section">
        <el-row :gutter="20">
          <el-col :span="isMobile ? 24 : 8">
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
              :class="{ 'mobile-input': isMobile }"
            />
          </el-col>
          <el-col :span="isMobile ? 24 : 4" v-if="!isMobile">
            <el-button @click="resetSearch">重置</el-button>
          </el-col>
        </el-row>
        <el-row v-if="isMobile" :gutter="10" style="margin-top: 10px;">
          <el-col :span="12">
            <el-button @click="resetSearch" class="mobile-btn">重置</el-button>
          </el-col>
          <el-col :span="12">
            <el-button type="primary" @click="handleSearch" class="mobile-btn">搜索</el-button>
          </el-col>
        </el-row>
      </div>

      <!-- 移动端用餐记录列表 -->
      <div v-if="isMobile" class="mobile-records-list" v-loading="mealRecordsStore.loading">
        <div
          v-for="record in mealRecordsStore.mealRecords"
          :key="record.id"
          class="mobile-record-card"
          @click="showRecordDetail(record)"
        >
          <div class="mobile-record-header">
            <div class="mobile-record-date">
              <el-icon><Calendar /></el-icon>
              {{ formatDate(record.created_at) }}
            </div>
            <div class="mobile-record-price">
              ¥{{ record.total_price }}
            </div>
          </div>

          <div class="mobile-record-image" v-if="record.image_url">
            <el-image
              :src="record.image_url"
              fit="cover"
              style="width: 100%; height: 120px; border-radius: 8px;"
              :preview-src-list="[record.image_url]"
            >
              <template #error>
                <div class="mobile-image-placeholder">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
          </div>

          <div class="mobile-record-content">
            <div class="mobile-dishes-preview">
              <span class="mobile-dishes-count">{{ record.dishes.length }}道菜</span>
              <div class="mobile-dishes-tags">
                <el-tag
                  v-for="(dishItem, index) in record.dishes.slice(0, 3)"
                  :key="dishItem.id"
                  size="small"
                  class="mobile-dish-tag"
                >
                  {{ dishItem.dish.name }}
                </el-tag>
                <el-tag v-if="record.dishes.length > 3" size="small" type="info">
                  +{{ record.dishes.length - 3 }}
                </el-tag>
              </div>
            </div>

            <div class="mobile-thoughts" v-if="record.thoughts">
              <p>{{ record.thoughts }}</p>
            </div>

            <div class="mobile-record-actions">
              <el-button 
                size="small" 
                type="primary"
                @click.stop="editRecord(record)"
              >
                编辑
              </el-button>
              <el-button 
                size="small" 
                type="danger"
                @click.stop="deleteRecord(record.id)"
              >
                删除
              </el-button>
            </div>
          </div>
        </div>

        <!-- 移动端加载更多 -->
        <div v-if="mealRecordsStore.hasMore" class="mobile-load-more">
          <el-button 
            type="primary" 
            :loading="mealRecordsStore.loading"
            @click="loadMore"
            class="mobile-btn"
          >
            加载更多
          </el-button>
        </div>
      </div>

      <!-- 桌面端用餐记录列表 -->
      <div v-else class="records-grid" v-loading="mealRecordsStore.loading">
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
            <div class="record-actions">
              <el-button
                size="small"
                type="primary"
                @click="editRecord(record)"
              >
                编辑
              </el-button>
              <el-button
                size="small"
                type="danger"
                @click="deleteRecord(record.id)"
              >
                删除
              </el-button>
            </div>
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

      <!-- 桌面端分页 -->
      <div class="pagination-wrapper" v-if="!isMobile">
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
    </div>

    <!-- 创建用餐记录对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      title="创建用餐记录"
      :width="isMobile ? '90%' : '600px'"
    >
      <el-form
        ref="recordFormRef"
        :model="recordForm"
        :rules="recordRules"
        label-width="100px"
      >
        <el-form-item label="选择菜品" prop="dish_ids">
          <el-checkbox-group v-model="recordForm.dish_ids" :class="{ 'mobile-checkbox-group': isMobile }">
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
      :width="isMobile ? '90%' : '800px'"
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
              class="dish-detail-item"
              @click="showDishDetail(dishItem.dish)"
            >
              <div class="dish-detail-image">
                <el-image
                  :src="dishItem.dish.image_url || '/placeholder-dish.jpg'"
                  fit="cover"
                  style="width: 100%; height: 120px; border-radius: 8px;"
                />
              </div>
              <div class="dish-detail-info">
                <h4>{{ dishItem.dish.name }}</h4>
                <p class="dish-detail-description">{{ dishItem.dish.description }}</p>
                <div class="dish-detail-meta">
                  <span class="dish-detail-price">¥{{ dishItem.dish.price }}</span>
                  <span class="dish-detail-category" v-if="dishItem.dish.category">
                    {{ dishItem.dish.category.name }}
                  </span>
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
          <h3>用餐照片</h3>
          <el-image
            :src="selectedRecord.image_url"
            fit="cover"
            style="width: 100%; max-height: 300px; border-radius: 8px;"
            :preview-src-list="[selectedRecord.image_url]"
          />
        </div>
      </div>
    </el-dialog>

    <!-- 编辑用餐记录对话框 -->
    <el-dialog
      v-model="showEditDialog"
      title="编辑用餐记录"
      :width="isMobile ? '90%' : '600px'"
    >
      <el-form
        ref="editFormRef"
        :model="editForm"
        :rules="recordRules"
        label-width="100px"
      >
        <el-form-item label="用餐感想" prop="thoughts">
          <el-input
            v-model="editForm.thoughts"
            type="textarea"
            :rows="3"
            placeholder="记录一下今天的用餐感受..."
          />
        </el-form-item>
        
        <el-form-item label="图片链接" prop="image_url">
          <el-input
            v-model="editForm.image_url"
            placeholder="可选：添加用餐图片链接"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showEditDialog = false">取消</el-button>
          <el-button type="primary" @click="updateRecord" :loading="updating">
            更新记录
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 菜品详情对话框 -->
    <el-dialog
      v-model="showDishDetailDialog"
      title="菜品详情"
      :width="isMobile ? '90%' : '600px'"
    >
      <div v-if="selectedDish" class="dish-detail">
        <div class="dish-detail-image">
          <el-image
            :src="selectedDish.image_url || '/placeholder-dish.jpg'"
            fit="cover"
            style="width: 100%; height: 200px;"
          />
        </div>
        <div class="dish-detail-info">
          <h3>{{ selectedDish.name }}</h3>
          <p class="description">{{ selectedDish.description }}</p>
          <div class="meta">
            <p><strong>价格:</strong> ¥{{ selectedDish.price }}</p>
            <p><strong>分类:</strong> {{ selectedDish.category?.name || '未分类' }}</p>
            <p v-if="selectedDish.cooking_link">
              <strong>制作链接:</strong> 
              <a :href="selectedDish.cooking_link" target="_blank">查看制作方法</a>
            </p>
          </div>
          <div v-if="selectedDish.ingredients && selectedDish.ingredients.length > 0" class="ingredients">
            <h4>所需食材:</h4>
            <el-tag
              v-for="ingredient in selectedDish.ingredients"
              :key="ingredient.id"
              class="ingredient-tag"
            >
              {{ ingredient.ingredient.name }} {{ ingredient.quantity }}{{ ingredient.ingredient.unit }}
            </el-tag>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, reactive } from 'vue'
import { useMealRecordsStore } from '@/stores/mealRecords'
import { useDishesStore } from '@/stores/dishes'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Calendar, Picture, ArrowLeft } from '@element-plus/icons-vue'

const mealRecordsStore = useMealRecordsStore()
const dishesStore = useDishesStore()

// 响应式数据
const dateRange = ref([])
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const showEditDialog = ref(false)
const showDishDetailDialog = ref(false)
const selectedRecord = ref(null)
const selectedDish = ref(null)
const editingRecord = ref(null)
const saving = ref(false)
const updating = ref(false)
const recordFormRef = ref()
const editFormRef = ref()
const isMobile = ref(false)

const recordForm = reactive({
  dish_ids: [],
  thoughts: '',
  image_url: ''
})

const editForm = reactive({
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

const editRecord = (record) => {
  editingRecord.value = record
  Object.assign(editForm, {
    thoughts: record.thoughts || '',
    image_url: record.image_url || ''
  })
  showEditDialog.value = true
}

const updateRecord = async () => {
  if (!editingRecord.value) return

  updating.value = true
  try {
    await mealRecordsStore.updateMealRecord(editingRecord.value.id, editForm)
    showEditDialog.value = false
    editingRecord.value = null
    ElMessage.success('用餐记录更新成功')
  } catch (error) {
    console.error('更新用餐记录失败:', error)
  } finally {
    updating.value = false
  }
}

// 加载更多
const loadMore = async () => {
  await mealRecordsStore.loadMore()
}

// 移动端检测
const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

// 生命周期
onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  mealRecordsStore.getMealRecords()
  dishesStore.getDishes()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
.meal-records-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

/* 桌面端样式 */
.page-header {
  background: white;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-header h1 {
  margin: 0;
  color: #333;
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.search-section {
  background: white;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.records-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
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
  gap: 8px;
  color: #666;
  font-size: 14px;
}

.record-actions {
  display: flex;
  gap: 10px;
}

.record-image {
  margin-bottom: 15px;
  border-radius: 8px;
  overflow: hidden;
}

.image-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  background-color: #f5f5f5;
  color: #999;
  gap: 10px;
}

.record-content {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.dishes-list h4 {
  margin: 0 0 10px 0;
  color: #333;
}

.dish-items {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.dish-tag {
  cursor: pointer;
  transition: background-color 0.3s;
}

.dish-tag:hover {
  background-color: #409EFF;
  color: white;
}

.total-price {
  font-size: 16px;
  color: #f56c6c;
}

.thoughts h4 {
  margin: 0 0 8px 0;
  color: #333;
}

.thoughts p {
  margin: 0;
  color: #666;
  line-height: 1.5;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}

/* 移动端样式 */
.mobile-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  padding: 10px 0;
}

.mobile-header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 15px;
}

.back-btn {
  font-size: 18px;
}

.mobile-header-content h2 {
  margin: 0;
  color: #333;
  font-size: 18px;
}

.mobile-content {
  padding-top: 60px;
  padding-bottom: 20px;
}

.mobile-records-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.mobile-record-card {
  background: white;
  border-radius: 12px;
  padding: 15px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
}

.mobile-record-card:active {
  transform: scale(0.98);
}

.mobile-record-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.mobile-record-date {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #666;
  font-size: 14px;
}

.mobile-record-price {
  color: #f56c6c;
  font-weight: bold;
  font-size: 16px;
}

.mobile-record-image {
  margin-bottom: 12px;
}

.mobile-image-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 120px;
  background-color: #f5f5f5;
  color: #999;
  font-size: 20px;
  border-radius: 8px;
}

.mobile-record-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.mobile-dishes-preview {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.mobile-dishes-count {
  color: #666;
  font-size: 12px;
}

.mobile-dishes-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.mobile-dish-tag {
  font-size: 10px;
}

.mobile-thoughts p {
  margin: 0;
  color: #666;
  font-size: 14px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.mobile-record-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

.mobile-load-more {
  text-align: center;
  margin: 20px 0;
}

/* 对话框样式 */
.record-detail {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 15px;
  border-bottom: 1px solid #eee;
}

.detail-date {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #666;
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
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 15px;
}

.dish-detail-item {
  background: #f9f9f9;
  border-radius: 8px;
  padding: 12px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.dish-detail-item:hover {
  background-color: #f0f0f0;
}

.dish-detail-image {
  margin-bottom: 10px;
}

.dish-detail-info h4 {
  margin: 0 0 6px 0;
  color: #333;
  font-size: 16px;
}

.dish-detail-description {
  color: #666;
  font-size: 12px;
  line-height: 1.3;
  margin-bottom: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.dish-detail-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dish-detail-price {
  color: #f56c6c;
  font-weight: bold;
}

.dish-detail-category {
  color: #409EFF;
  font-size: 10px;
  background-color: #ecf5ff;
  padding: 2px 6px;
  border-radius: 8px;
}

.detail-thoughts h3 {
  margin: 0 0 10px 0;
  color: #333;
}

.detail-thoughts p {
  margin: 0;
  color: #666;
  line-height: 1.6;
}

.detail-image h3 {
  margin: 0 0 10px 0;
  color: #333;
}

/* 菜品详情样式 */
.dish-detail {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.dish-detail-image {
  border-radius: 8px;
  overflow: hidden;
}

.dish-detail-info h3 {
  margin: 0 0 10px 0;
  color: #333;
}

.description {
  color: #666;
  line-height: 1.6;
  margin-bottom: 15px;
}

.meta p {
  margin: 8px 0;
  color: #333;
}

.meta a {
  color: #409EFF;
  text-decoration: none;
}

.ingredients {
  margin-top: 15px;
}

.ingredients h4 {
  margin: 0 0 10px 0;
  color: #333;
}

.ingredient-tag {
  margin: 0 5px 5px 0;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .main-content {
    padding: 10px;
  }
  
  .search-section {
    padding: 15px;
    margin-bottom: 15px;
  }
  
  .mobile-checkbox-group {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
  
  .mobile-checkbox-group .el-checkbox {
    margin-right: 0;
  }
  
  .mobile-input {
    width: 100%;
  }
  
  .mobile-btn {
    width: 100%;
    margin-bottom: 10px;
  }
}

@media (max-width: 480px) {
  .main-content {
    padding: 5px;
  }
  
  .mobile-record-card {
    padding: 12px;
  }
  
  .mobile-record-actions {
    flex-direction: column;
  }
}
</style> 