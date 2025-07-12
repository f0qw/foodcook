<template>
  <div class="home">
    <!-- 导航栏 -->
    <el-header class="header">
      <div class="header-content">
        <div class="logo">
          <h1>🍽️ FoodCook</h1>
        </div>
        <div class="nav">
          <el-menu mode="horizontal" :router="true" :default-active="$route.path">
            <el-menu-item index="/">首页</el-menu-item>
            <el-menu-item index="/dishes" v-if="authStore.isAuthenticated && authStore.isRoot">菜品管理</el-menu-item>
            <el-menu-item index="/ingredients" v-if="authStore.isAuthenticated && authStore.isRoot">食材管理</el-menu-item>
            <el-menu-item index="/meal-records" v-if="authStore.isAuthenticated">用餐记录</el-menu-item>
          </el-menu>
        </div>
        <div class="user-actions">
          <template v-if="authStore.isAuthenticated">
            <el-dropdown @command="handleCommand">
              <span class="user-info">
                {{ authStore.user?.username }}
                <el-icon><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                  <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <el-button type="primary" @click="$router.push('/login')">登录</el-button>
            <el-button @click="$router.push('/register')">注册</el-button>
          </template>
        </div>
      </div>
    </el-header>

    <!-- 主要内容 -->
    <div class="main-content">
      <!-- 搜索区域 -->
      <div class="search-section">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索菜品..."
          class="search-input"
          @keyup.enter="handleSearch"
        >
          <template #append>
            <el-button @click="handleSearch">
              <el-icon><Search /></el-icon>
            </el-button>
          </template>
        </el-input>
      </div>

      <!-- 菜品展示区域 -->
      <div class="dishes-section">
        <div class="section-header">
          <h2>今日推荐</h2>
          <div class="header-actions" v-if="authStore.isAuthenticated">
            <el-badge :value="cartStore.totalItems" :hidden="cartStore.totalItems === 0" class="cart-badge">
              <el-button type="info" @click="showCartDialog = true" :disabled="cartStore.totalItems === 0">
                <el-icon><ShoppingCart /></el-icon>
                购物车 ({{ cartStore.totalItems }})
              </el-button>
            </el-badge>
            <el-button type="primary" @click="showOrderDialog = true" :disabled="cartStore.totalItems === 0">
              创建用餐记录
            </el-button>
          </div>
        </div>

        <!-- 菜品网格 -->
        <el-row :gutter="20" v-loading="dishesStore.loading">
                      <el-col 
              v-for="dish in dishesStore.dishes" 
              :key="dish.id" 
              :xs="24" 
              :sm="12" 
              :md="8" 
              :lg="6"
              class="dish-col"
            >
              <el-card class="dish-card">
                <div class="dish-card-content" @click="showDishDetail(dish)">
                              <div class="dish-image">
                  <el-image
                    :src="dish.image_url || '/placeholder-dish.jpg'"
                    fit="cover"
                    :preview-src-list="[dish.image_url || '/placeholder-dish.jpg']"
                  >
                    <template #error>
                      <div class="image-placeholder">
                        <el-icon><Picture /></el-icon>
                      </div>
                    </template>
                  </el-image>
                </div>
                <div class="dish-info">
                  <h3>{{ dish.name }}</h3>
                  <p class="dish-description">{{ dish.description }}</p>
                  <div class="dish-meta">
                    <span class="price">¥{{ dish.price }}</span>
                    <span class="category" v-if="dish.category">{{ dish.category.name }}</span>
                  </div>
                </div>
                <div class="dish-actions">
                  <el-button 
                    v-if="authStore.isAuthenticated"
                    :type="cartStore.isInCart(dish.id) ? 'danger' : 'primary'"
                    size="small"
                    @click.stop="toggleCartItem(dish)"
                  >
                    <el-icon v-if="cartStore.isInCart(dish.id)"><Remove /></el-icon>
                    <el-icon v-else><Plus /></el-icon>
                    {{ cartStore.isInCart(dish.id) ? '移除' : '加入购物车' }}
                  </el-button>
                </div>
                </div>
              </el-card>
          </el-col>
        </el-row>

        <!-- 分页 -->
        <div class="pagination-wrapper">
          <el-pagination
            v-model:current-page="dishesStore.currentPage"
            v-model:page-size="dishesStore.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="dishesStore.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </div>

    <!-- 菜品详情对话框 -->
    <el-dialog
      v-model="dishDetailVisible"
      title="菜品详情"
      width="600px"
      :before-close="handleCloseDishDetail"
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

    <!-- 购物车对话框 -->
    <el-dialog
      v-model="showCartDialog"
      title="购物车"
      width="600px"
    >
      <div v-if="cartStore.totalItems === 0" class="empty-cart">
        <el-empty description="购物车是空的" />
      </div>
      <div v-else class="cart-content">
        <div class="cart-items">
          <div v-for="item in cartStore.cartItems" :key="item.id" class="cart-item">
            <div class="cart-item-info">
              <el-image
                :src="item.image_url || '/placeholder-dish.jpg'"
                style="width: 60px; height: 60px; border-radius: 8px;"
                fit="cover"
              />
              <div class="cart-item-details">
                <h4>{{ item.name }}</h4>
                <p class="cart-item-category" v-if="item.category">{{ item.category.name }}</p>
                <span class="cart-item-price">¥{{ item.price }}</span>
              </div>
            </div>
            <el-button 
              type="danger" 
              size="small" 
              @click="cartStore.removeFromCart(item.id)"
            >
              移除
            </el-button>
          </div>
        </div>
        <div class="cart-summary">
          <div class="cart-total">
            <span>总计: </span>
            <span class="total-price">¥{{ cartStore.totalPrice }}</span>
          </div>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCartDialog = false">关闭</el-button>
          <el-button type="primary" @click="createMealRecordFromCart" :loading="creating">
            创建用餐记录
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 创建用餐记录对话框 -->
    <el-dialog
      v-model="showOrderDialog"
      title="创建用餐记录"
      width="500px"
    >
      <div class="selected-dishes">
        <h4>已选择的菜品:</h4>
        <div v-for="item in cartStore.cartItems" :key="item.id" class="selected-dish">
          <span>{{ item.name }}</span>
          <span class="dish-price">¥{{ item.price }}</span>
        </div>
        <div class="total-section">
          <strong>总计: ¥{{ cartStore.totalPrice }}</strong>
        </div>
      </div>
      <el-form :model="orderForm" label-width="80px">
        <el-form-item label="用餐感想">
          <el-input
            v-model="orderForm.thoughts"
            type="textarea"
            :rows="3"
            placeholder="记录一下今天的用餐感受..."
          />
        </el-form-item>
        <el-form-item label="图片链接">
          <el-input
            v-model="orderForm.image_url"
            placeholder="可选：添加用餐图片链接"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showOrderDialog = false">取消</el-button>
          <el-button type="primary" @click="createMealRecord" :loading="creating">
            创建记录
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useDishesStore } from '@/stores/dishes'
import { useMealRecordsStore } from '@/stores/mealRecords'
import { useCartStore } from '@/stores/cart'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()
const dishesStore = useDishesStore()
const mealRecordsStore = useMealRecordsStore()
const cartStore = useCartStore()

// 响应式数据
const searchKeyword = ref('')
const dishDetailVisible = ref(false)
const selectedDish = ref(null)
const showCartDialog = ref(false)
const showOrderDialog = ref(false)
const creating = ref(false)

const orderForm = ref({
  dish_ids: [],
  thoughts: '',
  image_url: ''
})

// 计算属性
const totalPrice = computed(() => {
  return cartStore.totalPrice
})

// 方法
const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    dishesStore.searchDishes(searchKeyword.value.trim())
  } else {
    dishesStore.getDishes()
  }
}

const showDishDetail = (dish) => {
  selectedDish.value = dish
  dishDetailVisible.value = true
}

const handleCloseDishDetail = () => {
  dishDetailVisible.value = false
  selectedDish.value = null
}

const toggleCartItem = (dish) => {
  if (cartStore.isInCart(dish.id)) {
    cartStore.removeFromCart(dish.id)
    ElMessage.success('已从购物车移除')
  } else {
    cartStore.addToCart(dish)
    ElMessage.success('已添加到购物车')
  }
}

const handleSizeChange = (size) => {
  dishesStore.setPageSize(size)
  dishesStore.getDishes()
}

const handleCurrentChange = (page) => {
  dishesStore.setPage(page)
  dishesStore.getDishes()
}

const handleCommand = (command) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'logout':
      authStore.logout()
      router.push('/')
      break
  }
}

const createMealRecord = async () => {
  if (cartStore.totalItems === 0) {
    ElMessage.warning('请至少选择一个菜品')
    return
  }

  creating.value = true
  try {
    const mealRecordData = {
      dish_ids: cartStore.dishIds,
      thoughts: orderForm.value.thoughts,
      image_url: orderForm.value.image_url
    }
    await mealRecordsStore.createMealRecord(mealRecordData)
    showOrderDialog.value = false
    cartStore.clearCart()
    orderForm.value = {
      dish_ids: [],
      thoughts: '',
      image_url: ''
    }
    ElMessage.success(`用餐记录创建成功！总价格：¥${totalPrice.value}`)
  } catch (error) {
    console.error('创建用餐记录失败:', error)
  } finally {
    creating.value = false
  }
}

const createMealRecordFromCart = () => {
  showCartDialog.value = false
  showOrderDialog.value = true
}

// 生命周期
onMounted(() => {
  dishesStore.getDishes()
})
</script>

<style scoped>
.home {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.header {
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
}

.logo h1 {
  margin: 0;
  color: #409eff;
  font-size: 24px;
}

.nav {
  flex: 1;
  margin: 0 40px;
}

.user-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 5px;
  cursor: pointer;
  padding: 5px 10px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f5f5;
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.search-section {
  margin-bottom: 30px;
  text-align: center;
}

.search-input {
  max-width: 500px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h2 {
  margin: 0;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.cart-badge {
  margin-right: 10px;
}

.dish-col {
  margin-bottom: 20px;
}

.dish-card {
  transition: transform 0.3s, box-shadow 0.3s;
  height: 100%;
}

.dish-card-content {
  cursor: pointer;
}

.dish-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.dish-actions {
  padding: 10px 15px;
  border-top: 1px solid #f0f0f0;
  text-align: center;
}

.dish-image {
  height: 200px;
  overflow: hidden;
}

.dish-image .el-image {
  width: 100%;
  height: 100%;
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
  color: #999;
  font-size: 40px;
}

.dish-info {
  padding: 15px;
}

.dish-info h3 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 18px;
}

.dish-description {
  color: #666;
  font-size: 14px;
  margin-bottom: 15px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.dish-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.price {
  color: #f56c6c;
  font-size: 18px;
  font-weight: bold;
}

.category {
  background-color: #409eff;
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
}

.pagination-wrapper {
  text-align: center;
  margin-top: 30px;
}

.dish-detail {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.dish-detail-info h3 {
  margin: 0 0 10px 0;
  color: #333;
}

.description {
  color: #666;
  margin-bottom: 15px;
  line-height: 1.6;
}

.meta p {
  margin: 5px 0;
  color: #333;
}

.meta a {
  color: #409eff;
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
  margin: 5px 5px 5px 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.empty-cart {
  text-align: center;
  padding: 40px 0;
}

.cart-content {
  max-height: 400px;
  overflow-y: auto;
}

.cart-items {
  margin-bottom: 20px;
}

.cart-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #f0f0f0;
}

.cart-item:last-child {
  border-bottom: none;
}

.cart-item-info {
  display: flex;
  align-items: center;
  gap: 15px;
  flex: 1;
}

.cart-item-details h4 {
  margin: 0 0 5px 0;
  color: #333;
}

.cart-item-category {
  margin: 0 0 5px 0;
  color: #666;
  font-size: 12px;
}

.cart-item-price {
  color: #f56c6c;
  font-weight: bold;
}

.cart-summary {
  border-top: 2px solid #f0f0f0;
  padding-top: 15px;
}

.cart-total {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
}

.total-price {
  color: #f56c6c;
  font-weight: bold;
  font-size: 18px;
}

.selected-dishes {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 8px;
}

.selected-dishes h4 {
  margin: 0 0 10px 0;
  color: #333;
}

.selected-dish {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 5px 0;
  border-bottom: 1px solid #e9ecef;
}

.selected-dish:last-child {
  border-bottom: none;
}

.dish-price {
  color: #f56c6c;
  font-weight: bold;
}

.total-section {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 2px solid #dee2e6;
  text-align: right;
}
</style> 