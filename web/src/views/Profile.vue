<template>
  <div class="profile-page">
    <div class="page-header">
      <h1>个人中心</h1>
    </div>

    <div class="profile-content">
      <el-row :gutter="30">
        <!-- 个人信息卡片 -->
        <el-col :span="8">
          <el-card class="profile-card">
            <div class="avatar-section">
              <el-avatar :size="100" :src="userInfo.avatar_url">
                <el-icon><User /></el-icon>
              </el-avatar>
              <h2>{{ userInfo.username }}</h2>
              <p class="user-email">{{ userInfo.email }}</p>
            </div>
            
            <div class="user-stats">
              <div class="stat-item">
                <div class="stat-number">{{ stats.totalDishes }}</div>
                <div class="stat-label">菜品数量</div>
              </div>
              <div class="stat-item">
                <div class="stat-number">{{ stats.totalRecords }}</div>
                <div class="stat-label">用餐记录</div>
              </div>
              <div class="stat-item">
                <div class="stat-number">{{ stats.totalIngredients }}</div>
                <div class="stat-label">食材数量</div>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 详细信息 -->
        <el-col :span="16">
          <el-card>
            <template #header>
              <div class="card-header">
                <span>个人信息</span>
                <el-button type="primary" @click="showEditDialog = true">
                  编辑信息
                </el-button>
              </div>
            </template>

            <el-descriptions :column="2" border>
              <el-descriptions-item label="用户名">
                {{ userInfo.username }}
              </el-descriptions-item>
              <el-descriptions-item label="邮箱">
                {{ userInfo.email }}
              </el-descriptions-item>
              <el-descriptions-item label="注册时间">
                {{ formatDate(userInfo.created_at) }}
              </el-descriptions-item>
              <el-descriptions-item label="最后登录">
                {{ formatDate(userInfo.updated_at) }}
              </el-descriptions-item>
            </el-descriptions>
          </el-card>

          <!-- 最近用餐记录 -->
          <el-card style="margin-top: 20px;">
            <template #header>
              <div class="card-header">
                <span>最近用餐记录</span>
                <el-button @click="$router.push('/meal-records')">查看全部</el-button>
              </div>
            </template>

            <div v-if="recentRecords.length > 0">
              <div
                v-for="record in recentRecords"
                :key="record.id"
                class="recent-record"
              >
                <div class="record-info">
                  <div class="record-date">{{ formatDate(record.created_at) }}</div>
                  <div class="record-dishes">
                    {{ record.dishes.map(d => d.name).join(', ') }}
                  </div>
                  <div class="record-price">¥{{ calculateTotal(record.dishes) }}</div>
                </div>
              </div>
            </div>
            <el-empty v-else description="暂无用餐记录" />
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 编辑个人信息对话框 -->
    <el-dialog
      v-model="showEditDialog"
      title="编辑个人信息"
      width="500px"
    >
      <el-form
        ref="editFormRef"
        :model="editForm"
        :rules="editRules"
        label-width="100px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="editForm.username" />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="editForm.email" />
        </el-form-item>
        
        <el-form-item label="头像链接" prop="avatar_url">
          <el-input v-model="editForm.avatar_url" placeholder="可选：输入头像图片链接" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showEditDialog = false">取消</el-button>
          <el-button type="primary" @click="saveProfile" :loading="saving">
            保存
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useDishesStore } from '@/stores/dishes'
import { useMealRecordsStore } from '@/stores/mealRecords'
import { ElMessage } from 'element-plus'

const authStore = useAuthStore()
const dishesStore = useDishesStore()
const mealRecordsStore = useMealRecordsStore()

// 响应式数据
const showEditDialog = ref(false)
const saving = ref(false)
const editFormRef = ref()

const userInfo = ref({
  username: '',
  email: '',
  avatar_url: '',
  created_at: '',
  updated_at: ''
})

const editForm = reactive({
  username: '',
  email: '',
  avatar_url: ''
})

const editRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

const recentRecords = ref([])

// 计算属性
const stats = computed(() => ({
  totalDishes: dishesStore.total,
  totalRecords: mealRecordsStore.total,
  totalIngredients: 0 // 这里可以添加食材统计
}))

// 方法
const loadUserInfo = async () => {
  try {
    const user = await authStore.getProfile()
    userInfo.value = user
    Object.assign(editForm, {
      username: user.username,
      email: user.email,
      avatar_url: user.avatar_url || ''
    })
  } catch (error) {
    console.error('获取用户信息失败:', error)
  }
}

const loadRecentRecords = async () => {
  try {
    const response = await mealRecordsStore.getMealRecords({ limit: 5 })
    recentRecords.value = response.data || []
  } catch (error) {
    console.error('获取最近用餐记录失败:', error)
  }
}

const saveProfile = async () => {
  if (!editFormRef.value) return

  try {
    await editFormRef.value.validate()
    saving.value = true

    // 这里应该调用更新用户信息的API
    // await authAPI.updateProfile(editForm)
    
    ElMessage.success('个人信息更新成功')
    showEditDialog.value = false
    await loadUserInfo()
  } catch (error) {
    console.error('更新个人信息失败:', error)
  } finally {
    saving.value = false
  }
}

const formatDate = (dateString) => {
  if (!dateString) return '未知'
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
onMounted(async () => {
  await loadUserInfo()
  await loadRecentRecords()
  await dishesStore.getDishes()
  await mealRecordsStore.getMealRecords()
})
</script>

<style scoped>
.profile-page {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h1 {
  margin: 0;
  color: #333;
}

.profile-content {
  margin-top: 20px;
}

.profile-card {
  text-align: center;
}

.avatar-section {
  margin-bottom: 30px;
}

.avatar-section h2 {
  margin: 15px 0 5px 0;
  color: #333;
}

.user-email {
  margin: 0;
  color: #666;
  font-size: 14px;
}

.user-stats {
  display: flex;
  justify-content: space-around;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.stat-item {
  text-align: center;
}

.stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 12px;
  color: #666;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.recent-record {
  padding: 15px 0;
  border-bottom: 1px solid #f0f0f0;
}

.recent-record:last-child {
  border-bottom: none;
}

.record-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.record-date {
  color: #666;
  font-size: 14px;
  width: 150px;
}

.record-dishes {
  flex: 1;
  margin: 0 20px;
  color: #333;
}

.record-price {
  color: #f56c6c;
  font-weight: bold;
  width: 80px;
  text-align: right;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 